package sd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client"
	"log"
	"strings"
	"sync"
	"time"
)

const service = "service_discovery"

type Master struct {
	sync.RWMutex
	kapi  client.KeysAPI
	key   string
	nodes map[string]string
}

func NewMaster(endpoints []string) *Master {
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
		// set timeout for per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	master := &Master{
		kapi:  client.NewKeysAPI(c),
		key:   fmt.Sprintf("/%s", service),
		nodes: make(map[string]string),
	}

	err = master.InitNodes()
	if err != nil {
		log.Fatal(err)
	}

	go master.WatchWorkers()

	return master
}

func (m *Master) InitNodes() error {
	resp, err := m.kapi.Get(context.Background(), m.key, nil)
	if err != nil {
		return err
	}

	if resp.Node.Dir {
		for _, n := range resp.Node.Nodes {
			m.nodes[n.Key] = n.Value
		}
	}

	log.Println("Init nodes succeed, nodes:", m.nodes)

	return nil
}

func (m *Master) WatchWorkers() {
	watcher := m.kapi.Watcher(service, &client.WatcherOptions{
		Recursive: true,
	})

	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("Error watch workers: ", err)
			continue
		}

		// fmt.Println("action:", resp.Action, "key:", resp.Node.Key, "value:", resp.Node.Value)

		switch resp.Action {
		case "set", "update":
			m.addNode(resp.Node.Key, resp.Node.Value)
			break
		case "expire", "delete":
			m.deleteNode(resp.Node.Key)
			break
		default:
			log.Printf("Master watches unknown resp, key: %s, value: %s", resp.Node.Key, resp.Node.Value)
		}
	}
}

func (m *Master) GetNodes() (map[string]string, error) {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.nodes, nil
}

func (m *Master) addNode(key string, value string) {
	m.Lock()
	defer m.Unlock()
	// note: strings.TrimLeft("/sd/node1", "/sd") returns "node1"
	key = strings.TrimLeft(key, m.key)
	m.nodes[key] = value
}

func (m *Master) deleteNode(key string) {
	m.Lock()
	defer m.Unlock()
	key = strings.TrimLeft(key, m.key)
	delete(m.nodes, key)
}

package sd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client"
	"log"
	"sync"
	"time"
)

type Worker struct {
	sync.Mutex
	name              string
	kapi              client.KeysAPI
	key               string
	value             string
	heartbeatInterval time.Duration
	active            bool
}

func NewWorker(name string, key string, value string, endpoints []string) (*Worker, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: time.Second * 2,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	worker := &Worker{
		name:              name,
		kapi:              client.NewKeysAPI(c),
		key:               fmt.Sprintf("/%s/%s", service, key),
		value:             value,
		heartbeatInterval: time.Second * 2,
		active:            false,
	}

	return worker, nil
}

func (w *Worker) Register() {
	w.active = true
	w.periodHeartbeat()
}

func (w *Worker) UnRegister() {
	w.Lock()
	defer w.Unlock()
	w.active = false
}

func (w *Worker) periodHeartbeat() {
	for w.active {
		w.heartbeat()
		time.Sleep(w.heartbeatInterval)
	}
}

func (w *Worker) heartbeat() {
	w.Lock()
	defer w.Unlock()
	_, err := w.kapi.Set(context.Background(), service+w.key, w.value, &client.SetOptions{
		TTL: time.Second * 10, // TTL defines a period of time after-which the Node should expire and no longer exist
	})
	if err != nil {
		log.Println("Worker send heartbeat failed:", err)
	}
}

func (w *Worker) GetName() string {
	return w.name
}

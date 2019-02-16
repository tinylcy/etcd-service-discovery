package main

import (
	"etcd-service-discovery/sd"
	"log"
	"time"
)

func main() {
	master := sd.NewMaster([]string{"http://127.0.0.1:2379"})
	go master.WatchWorkers()

	for {
		nodes, err := master.GetNodes()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Master is running, nodes: %v", nodes)
		time.Sleep(time.Second * 2)
	}
}

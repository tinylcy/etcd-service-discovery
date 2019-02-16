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
		log.Println("...")
		time.Sleep(time.Second * 2)
	}
}

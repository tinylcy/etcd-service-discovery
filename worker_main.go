package main

import (
	"etcd-service-discovery/sd"
	"log"
	"time"
)

func main() {
	worker, err := sd.NewWorker("worker1", "node100", "v1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	go worker.Register()

	for {
		log.Printf("Worker [%s] is running...", worker.GetName())
		time.Sleep(time.Second * 2)
	}
}

package client

import (
	"github.com/OpenIoTHub/gateway-go/v2/tasks"
	"log"
)

var IsLibrary = true

func Run() {
	go start()
}

func start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("gateway-go panic:%+v", err)
		}
	}()
	tasks.RunTasks()
	go startHTTP()
	startGRPC()
}

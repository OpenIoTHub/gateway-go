package main

import (
	"github.com/OpenIoTHub/gateway-go/v2/client"
	"time"
)

func main() {
	client.Run()
	time.Sleep(500 * time.Second)
}

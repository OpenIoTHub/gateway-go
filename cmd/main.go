package main

import (
	"time"

	"github.com/OpenIoTHub/gateway-go/v2/client"
)

func main() {
	client.Run()
	time.Sleep(500 * time.Second)
}

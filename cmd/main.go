package main

import (
	"github.com/OpenIoTHub/gateway-go/v2/client"
	"github.com/OpenIoTHub/gateway-go/v2/register"
	"net"
	"time"
)

func main() {
	client.Run()
	register.RegisterService("localhost-gateway-go",
		"_http._tcp",
		"local",
		"localhost",
		client.HttpPort,
		[]string{"name=gateway-go", "home-page=https://github.com/OpenIoTHub/gateway-go"},
		0,
		[]net.IP{net.ParseIP("127.0.0.1")},
		[]net.IP{},
	)
	time.Sleep(500 * time.Second)
}

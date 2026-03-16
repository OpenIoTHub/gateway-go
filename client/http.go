package client

import (
	"fmt"
	"log"
	"net"

	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/register"
	"github.com/OpenIoTHub/gateway-go/v2/services"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func startHTTP() {
	port := config.ConfigMode.HttpServicePort
	register.RegisterService("localhost-gateway-go",
		"_http._tcp",
		"local",
		"localhost",
		port,
		[]string{
			"name=gateway-go",
			fmt.Sprintf("id=gateway-go@%s", uuid.Must(uuid.NewV4()).String()),
			"home-page=https://github.com/OpenIoTHub/gateway-go",
		},
		0,
		[]net.IP{net.ParseIP("127.0.0.1")},
		[]net.IP{},
	)
	r := gin.Default()
	r.GET("/", services.GatewayManager.IndexHandler)
	r.GET("/DisplayQrHandler", services.GatewayManager.DisplayQrHandler)
	log.Printf("Http 监听端口: %d\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Printf("HTTP 服务启动失败: %v", err)
	}
}

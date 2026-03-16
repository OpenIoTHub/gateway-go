package client

import (
	"fmt"
	"log"
	"net"

	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/info"
	"github.com/OpenIoTHub/gateway-go/v2/services"
	"github.com/grandcat/zeroconf"
)

func registerGatewayMDNS(gRpcPort int) {
	mac := "mac"
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	} else if len(interfaces) > 0 {
		mac = interfaces[0].HardwareAddr.String()
	}

	gatewayUUID, serverHost, _ := services.GatewayManager.GetLoginInfo()

	_, err = zeroconf.Register(
		fmt.Sprintf("OpenIoTHubGateway-%s", config.ConfigMode.GatewayUUID),
		"_openiothub-gateway._tcp",
		"local.",
		gRpcPort,
		[]string{
			"name=网关",
			"model=com.iotserv.services.gateway",
			fmt.Sprintf("mac=%s", mac),
			fmt.Sprintf("id=%s", config.ConfigMode.GatewayUUID),
			fmt.Sprintf("run_id=%s", gatewayUUID),
			fmt.Sprintf("server_host=%s", serverHost),
			"author=Farry",
			"email=newfarry@126.com",
			"home-page=https://github.com/OpenIoTHub",
			"firmware-respository=https://github.com/OpenIoTHub/gateway-go/v2",
			fmt.Sprintf("firmware-version=%s", info.Version),
		},
		nil,
	)
	if err != nil {
		log.Printf("mDNS 注册失败: %v", err)
	}
}

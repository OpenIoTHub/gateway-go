package docker

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/moby/moby/client"
)

// GetContainersInfo 获取当前主机上所有Docker容器的信息
func GetContainersInfo() (client.ContainerListResult, error) {
	// 创建Docker客户端
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return client.ContainerListResult{}, fmt.Errorf("failed to create docker client: %v", err)
	}
	defer cli.Close()

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取容器列表（包括停止的容器）
	return cli.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
}

func GetContainersServices() (rst []models.MDNSResult) {
	clr, err := GetContainersInfo()
	if err != nil {
		return []models.MDNSResult{}
	}
	for _, items := range clr.Items {
		port := 0
		for _, portInfo := range items.Ports {
			if portInfo.Type == "tcp" {
				port = int(portInfo.PublicPort)
			}
		}
		name := "Docker Service"
		for _, n := range items.Names {
			name = strings.ReplaceAll(n, "/", "")
			break
		}
		rst = append(rst, models.MDNSResult{
			Instance: items.ID,
			Service:  "_http._tcp",
			Domain:   "local",
			HostName: "localhost",
			Port:     port,
			Text:     []string{fmt.Sprintf("name=%s", name), fmt.Sprintf("id=%s", items.ID)},
			TTL:      0,
			AddrIPv4: []net.IP{net.ParseIP("127.0.0.1")},
			AddrIPv6: []net.IP{},
		})
	}
	return
}

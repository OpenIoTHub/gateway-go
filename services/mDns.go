package services

import (
	"context"
	"encoding/json"
	"fmt"
	"git.iotserv.com/iotserv/utils/models"
	"git.iotserv.com/iotserv/utils/msg"
	"github.com/iotdevice/zeroconf"
	"net"
	"time"
)

func findAllmDNS(stream net.Conn, service *models.NewService) error {
	//decode json
	var config *models.FindmDNS
	var rst []*zeroconf.ServiceEntry
	err := json.Unmarshal([]byte(service.Config), &config)
	if err != nil {
		return err
	}
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Println(entry)
			rst = append(rst, entry)
		}
	}(entries)
	timeOut := time.Millisecond * time.Duration(config.Second) * 150
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err = resolver.Browse(ctx, config.Service, config.Domain, entries)
	if err != nil {
		return err
	}
	<-ctx.Done()
	//fmt.Println("获取完成：")
	//if len(rst) > 0 {
	//	fmt.Println(rst[0])
	//}
	rstByte, err := json.Marshal(&rst)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		fmt.Println("写消息错误：")
		fmt.Println(err.Error())
	}
	return err
}

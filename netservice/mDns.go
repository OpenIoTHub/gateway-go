package netservice

import (
	"context"
	"encoding/json"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"github.com/grandcat/zeroconf"
	"log"
	"net"
	"time"
)

func FindAllmDNS(stream net.Conn, service *models.NewService) error {
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
			log.Println(entry)
			//TODO 去掉记录中ip不是本网段的ip
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
	//log.Println("获取完成：")
	//if len(rst) > 0 {
	//	log.Println(rst[0])
	//}
	rstByte, err := json.Marshal(&rst)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println(string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	return err
}

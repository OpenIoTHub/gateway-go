package mdns

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/OpenIoTHub/gateway-go/v2/register"
	"github.com/OpenIoTHub/utils/v2/models"
	"github.com/OpenIoTHub/utils/v2/msg"
	"github.com/grandcat/zeroconf"
)

func (mc *MdnsCtrl) FindAllmDNS(stream net.Conn, service *models.NewService) error {
	//decode json
	var config *models.FindmDNS
	var rst = make([]*models.MDNSResult, 0)
	err := json.Unmarshal([]byte(service.Config), &config)
	if err != nil {
		log.Println("json.Unmarshal([]byte(service.Config), &config):" + err.Error())
		return err
	}

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Println("zeroconf.NewResolver:" + err.Error())
	} else {
		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				//log.Println("entry:", entry)
				//TODO 去掉记录中ip不是本网段的ip
				rst = append(rst, &models.MDNSResult{
					Instance: entry.Instance,
					Service:  entry.Service,
					Domain:   entry.Domain,
					HostName: entry.HostName,
					Port:     entry.Port,
					Text:     entry.Text,
					TTL:      entry.TTL,
					AddrIPv4: entry.AddrIPv4,
					AddrIPv6: entry.AddrIPv6,
				})
			}
		}(entries)
		//TODO 发现时间
		timeOut := time.Millisecond * time.Duration(config.Second) * 250
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		err = resolver.Browse(ctx, config.Service, config.Domain, entries)
		if err != nil {
			log.Println("resolver.Browse:" + err.Error())
			return err
		}
		<-ctx.Done()
	}
	registeredServices := register.GetRegisteredServices()
	//发现_services._dns-sd._udp类型的时候添加所有手动注册的类型
	if config.Service == "_services._dns-sd._udp" {
		registeredType := make([]string, 0)
	Loop1:
		for _, registeredService := range registeredServices {
			for _, item := range registeredType {
				if item == registeredService.Service {
					continue Loop1
				}
			}
			//log.Println("ADD Registered service: ", registeredService.Service)
			rst = append(rst, &models.MDNSResult{
				Instance: registeredService.Service + ".local",
				Service:  "_services._dns-sd._udp",
				Domain:   "local",
			})
			registeredType = append(registeredType, registeredService.Service)
		}
	} else {
		for _, registeredService := range registeredServices {
			rst = append(rst, &models.MDNSResult{
				Instance: registeredService.Instance,
				Service:  registeredService.Service,
				Domain:   registeredService.Domain,
				HostName: registeredService.HostName,
				Port:     registeredService.Port,
				Text:     registeredService.Text,
				TTL:      registeredService.TTL,
				AddrIPv4: registeredService.AddrIPv4,
				AddrIPv6: registeredService.AddrIPv6,
			})
		}
	}
	rstByte, err := json.Marshal(&rst)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	//log.Println("mdns rstByte:", string(rstByte))
	err = msg.WriteMsg(stream, &models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	if err != nil {
		log.Println("写消息错误：")
		log.Println(err.Error())
	}
	//content, _ := json.Marshal(&models.JsonResponse{Code: 0, Msg: "Success", Result: string(rstByte)})
	//log.Println("content:", string(content))
	//stream.Close()
	return err
}

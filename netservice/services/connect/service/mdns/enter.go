package mdns

import (
	"context"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/grandcat/zeroconf"
	"log"
	"strings"
	"sync"
	"time"
)

var MdnsManager *MdnsCtrl

// MdnsCtrl mdns管理器
type MdnsCtrl struct {
	sync.Mutex
	serviceTypeList       map[string]bool
	serviceTypeServiceMap map[string][]models.MDNSResult
}

func init() {
	//NewMdnsManager()
}

func NewMdnsManager() {
	if MdnsManager == nil {
		MdnsManager = new(MdnsCtrl)
		MdnsManager.serviceTypeList = make(map[string]bool)
		MdnsManager.serviceTypeServiceMap = make(map[string][]models.MDNSResult)
		MdnsManager.startService()
	}
}

func (mc *MdnsCtrl) startService() {
	go mc.startZeroconfFind()
	go mc.startAvahiFind()
}

func (mc *MdnsCtrl) startZeroconfFind() (err error) {
	var resolver *zeroconf.Resolver
	for {
		const mdnsFindTypes = "_services._dns-sd._udp"
		resolver, err = zeroconf.NewResolver(nil)
		if err != nil {
			return err
		}

		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				mc.serviceTypeList[strings.Replace(entry.Instance, ".local", "", 1)] = true
			}
		}(entries)

		//TODO 发现时间
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err = resolver.Browse(ctx, mdnsFindTypes, "local", entries)
		if err != nil {
			return err
		}
		<-ctx.Done()
		//close(entries)

		for key, _ := range mc.serviceTypeList {
			entries2 := make(chan *zeroconf.ServiceEntry)
			go func(results <-chan *zeroconf.ServiceEntry) {
				for entry := range results {
					//去重添加
					mc.serviceTypeServiceMap[entry.Service] = append(mc.serviceTypeServiceMap[entry.Service], models.MDNSResult{
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
			}(entries2)
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
			resolver, err = zeroconf.NewResolver(nil)
			if err != nil {
				log.Println(err)
				cancel()
				continue
			}
			err = resolver.Browse(ctx, key, "local", entries2)
			if err != nil {
				log.Println(err)
				cancel()
				continue
			}
			<-ctx.Done()
			cancel()
			time.Sleep(time.Second)
		}
		//close(entries2)
		fmt.Printf("mdns:%+v", mc.serviceTypeServiceMap)
	}
}

func (mc *MdnsCtrl) startAvahiFind() (err error) {
	return
}

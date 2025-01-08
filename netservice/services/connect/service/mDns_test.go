package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	"github.com/grandcat/zeroconf"
	mdns "github.com/hashicorp/mdns"
	"log"
	"testing"
	"time"
)

func TestFindAllmDNS(t *testing.T) {
	var rst = make([]*models.MDNSResult, 0)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		panic(err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Printf("entry:%+v", entry)
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
	timeOut := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	//err = resolver.Browse(ctx, "_services._dns-sd._udp", "local", entries)
	err = resolver.Browse(ctx, "_airplay._tcp", "local", entries)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
	//log.Println("获取完成：")
	//if len(rst) > 0 {
	//	log.Println(rst[0])
	//}
	rstByte, err := json.Marshal(&rst)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	fmt.Println(string(rstByte))
}

func TestFindAllmDNS2(t *testing.T) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %+v\n", *entry)
		}
	}()
	params := mdns.DefaultParams("_airplay._tcp")
	params.Entries = entriesCh
	params.Timeout = 5 * time.Second
	params.DisableIPv6 = true
	// Start the lookup
	//err := mdns.Lookup("_airplay._tcp", entriesCh)
	err := mdns.Query(params)
	if err != nil {
		log.Println(err.Error())
	}
	time.Sleep(1 * time.Second)

	close(entriesCh)
}

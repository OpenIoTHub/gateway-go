package component

import (
	"context"
	"github.com/grandcat/zeroconf"
	"log"
	"time"
)

func CheckComponentExist(instance string) (bool, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
		return false, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	//TODO 是否需要手动关闭channel？
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()
	err = resolver.Browse(ctx, "_iotdevice._tcp", "local", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
		return false, err
	}

	<-ctx.Done()
	for entry := range entries {
		if entry.Instance == instance {
			return true, nil
		}
	}
	return false, nil
}

package component

import (
	"context"
	"github.com/grandcat/zeroconf"
	"log"
	"strings"
	"time"
)

func CheckComponentExist(model string) (bool, error) {
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
		for _, text := range entry.Text {
			keyValue := strings.Split(text, "=")
			if len(keyValue) == 2 && keyValue[0] == "model" && keyValue[1] == model {
				return true, nil
			}
		}
	}
	return false, nil
}

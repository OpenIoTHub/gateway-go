//mqtt消息组件
package component

import (
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/surgemq/surgemq/service"
)

func init() {
	svr := &service.Server{
		KeepAlive:        3600,
		ConnectTimeout:   30,
		AckTimeout:       30,
		TimeoutRetries:   5,
		SessionsProvider: "mem",
		TopicsProvider:   "mem",
	}
	port := 1884
	mqttaddr := fmt.Sprintf("tcp://0.0.0.0:%d", port)
	go func() {
		err := svr.ListenAndServe(mqttaddr)
		if err != nil {
			fmt.Printf("surgemq: %v", err)
		}
	}()
	_, err := zeroconf.Register("mqtt", "_component._tcp", "local.", port, []string{}, nil)
	if err != nil {
		fmt.Printf("zeroconf: %v", err)
	}
}

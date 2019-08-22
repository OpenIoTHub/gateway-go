//mqtt消息组件
package component

import (
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/surgemq/surgemq/service"
)

func init() {
	var componentName = "mqtt"
	//TODO 判断本网络是否已经存在此类型的组件，存在则不启动
	exist, err := CheckComponentExist(componentName)
	if err != nil {
		fmt.Println(err.Error())
	}
	if exist {
		fmt.Printf("本网络已经存在了%s组件，不启动/n", componentName)
	}
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
	server, err := zeroconf.Register(componentName, "_component._tcp", "local.", port, []string{}, nil)
	if err != nil {
		fmt.Printf("zeroconf: %v", err)
		return
	}
	err = svr.ListenAndServe(mqttaddr)
	if err != nil {
		fmt.Printf("surgemq: %v", err)
		server.Shutdown()
		return
	}
}

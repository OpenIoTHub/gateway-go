//mqtt消息组件
package component

import (
	"fmt"
	"github.com/OpenIoTHub/surgemq/service"
	nettool "github.com/OpenIoTHub/utils/net"
	"log"
	"runtime"
)

func init() {
	if runtime.GOOS == "android" {
		return
	}
	go func() {
		var txtInfo = nettool.MDNSServiceBaseInfo
		var model = "com.iotserv.services.mqttd"
		port, err := nettool.GetOneFreeTcpPort()
		if err != nil {
			log.Println(err.Error())
			return
		}
		//TODO 判断本网络是否已经存在此类型的组件，存在则不启动
		exist, err := nettool.CheckComponentExist(model)
		if err != nil {
			log.Println(err.Error())
		}
		if exist {
			fmt.Printf("本网络已经存在了%s组件，不启动\n", model)
			return
		}
		svr := &service.Server{
			KeepAlive:        3600,
			ConnectTimeout:   30,
			AckTimeout:       30,
			TimeoutRetries:   5,
			SessionsProvider: "mem",
			TopicsProvider:   "mem",
		}
		mqttaddr := fmt.Sprintf("tcp://0.0.0.0:%d", port)
		txtInfo["name"] = "mqtt服务器"
		txtInfo["model"] = model
		server, err := nettool.RegistermDNSService(txtInfo, port)
		err = svr.ListenAndServe(mqttaddr)
		if err != nil {
			fmt.Printf("surgemq: %v", err)
			server.Shutdown()
			return
		}
	}()
}

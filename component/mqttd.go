//mqtt消息组件
package component

import (
	"fmt"
	nettool "git.iotserv.com/iotserv/utils/net"
	"github.com/grandcat/zeroconf"
	"github.com/satori/go.uuid"
	"github.com/surgemq/surgemq/service"
)

func init() {
	go func() {
		var model = "com.iotserv.devices.mqttd"
		var port = 1884
		//TODO 判断本网络是否已经存在此类型的组件，存在则不启动
		exist, err := CheckComponentExist(model)
		if err != nil {
			fmt.Println(err.Error())
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
		var mac = uuid.Must(uuid.NewV4()).String()
		fmt.Println("mac:", mac)
		macs, err := nettool.GetMacs()
		fmt.Println("macs:", macs)
		if err == nil && len(macs) > 0 {
			fmt.Println("mac-len:", len(macs))
			for _, vMac := range macs {
				if vMac != "" {
					mac = vMac
				}
			}
		}
		var txt = []string{
			"name=mqtt服务器",
			fmt.Sprintf("model=%s", model),
			fmt.Sprintf("mac=%s", mac),
			fmt.Sprintf("id=%s", uuid.Must(uuid.NewV4()).String()),
			//web,native,none
			"ui-support=none",
			"ui-first=none",
			"author=Farry",
			"email=newfarry@126.com",
			"home-page=https://github.com/iotdevice",
			"firmware-respository=https://github.com/OpenIoTHub/GateWay",
			"firmware-version=1.0",
		}
		server, err := zeroconf.Register(fmt.Sprintf("%s-%s", model, mac), "_iotdevice._tcp", "local.", port, txt, nil)
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
	}()
}

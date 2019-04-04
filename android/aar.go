package android

import (
	"fmt"
	"git.iotserv.com/iotserv/client/config"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"github.com/satori/go.uuid"
	"net/http"
)

var tkstr = ""
var tkstr2 = ""

func token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, tkstr2)
}

func Run() {
	//crypto.Salt = "abc"
	id := uuid.Must(uuid.NewV4()).String()
	host := config.RegisterHost
	tkstr, _ = crypto.GetToken(id, host, config.TcpPort, config.KcpPort, config.TlsPort, config.UdpApiPort, 1, 200000000000) //47.96.185.226，118.89.106.226
	_, err := services.RunNATManager(tkstr)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	} else {
		fmt.Printf("登陆成功！\n")
	}
	//fmt.Printf(tkstr+"\n")
	tkstr2, _ = crypto.GetToken(id, host, config.TcpPort, config.KcpPort, config.TlsPort, config.UdpApiPort, 2, 200000000000) //118.89.106.226
	fmt.Printf("要想访问本内网，请用explorer使用以下token：\n" + tkstr2 + "\n")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.HandleFunc("/", token) //设置访问的路由

	err = http.ListenAndServe("127.0.0.1:1082", nil) //设置监听的端口
	if err != nil {
		fmt.Printf("请检查端口1082是否被占用")
	}
}

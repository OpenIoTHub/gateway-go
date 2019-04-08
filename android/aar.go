package android

import (
	"fmt"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/satori/go.uuid"
	"net/http"
)

var clientToken = ""
var explorerToken = ""
var configMode = models.ClientConfig{}

func init() {
	configMode.LastId = uuid.Must(uuid.NewV4()).String()
	configMode.Server = models.Srever{
		"tcp",
		"s1.365hour.com",
		34320,
		34320,
		34321,
		34321,
		"HLLdsa544&*S",
	}
}

func Run() {
	clientToken, err := crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 1, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = services.RunNATManager(configMode.Server.LoginKey, clientToken)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	} else {
		fmt.Printf("登陆成功！\n")
	}
	//fmt.Printf(tkstr+"\n")
	explorerToken, err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 2, 200000000000)
	fmt.Printf("要想访问本内网，请用explorer使用以下token：\n" + explorerToken + "\n")
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

func token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, explorerToken)
}

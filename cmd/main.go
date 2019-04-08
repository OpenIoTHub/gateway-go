package main

import (
	"flag"
	"fmt"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	h              bool
	configFilePath string
	p              string
)
var clientToken = ""
var explorerToken = ""

func init() {
	flag.BoolVar(&h, "h", false, "show this help")
	flag.StringVar(&configFilePath, "c", "./client.yaml", "set `config`")
	flag.StringVar(&p, "p", "1082", "set `port`")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `nat-cloud.com 内网管理端，运行在需要被穿透的内网的一个主机上。

Usage: nat -h
		nat -host
		nat -t token [-p port] 

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}
	configMode := models.ClientConfig{}
	_, err := os.Stat(configFilePath)
	if err != nil {
		fmt.Println("没有找到配置文件：", configFilePath)
		fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
		//	生成配置文件模板
		configMode.LastId = uuid.Must(uuid.NewV4()).String()
		configMode.LastExplorerToken = "不需要配置此项，随后自动生成"
		err = writeConfigFile(configMode, configFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	}
	//配置文件存在
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//找到了配置文件
	if len(configMode.LastId) < 35 {
		configMode.LastId = uuid.Must(uuid.NewV4()).String()
	}
	clientToken, err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 1, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = services.RunNATManager(configMode.Server.LoginKey, clientToken)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		os.Exit(0)
	}
	fmt.Printf("登陆成功！\n")
	explorerToken, err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 2, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("访问token：\n\n" + explorerToken + "\n\n")
	configMode.LastExplorerToken = explorerToken
	err = writeConfigFile(configMode, configFilePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("你也可以访问：http://127.0.0.1:%s/查看访问token\n", p)
	http.HandleFunc("/", token)                  //设置访问的路由
	err = http.ListenAndServe("0.0.0.0:"+p, nil) //设置监听的端口
	if err != nil {
		fmt.Printf("请检查端口" + p + "是否被占用")
	}
}

func writeConfigFile(configMode models.ClientConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		fmt.Println("写入配置文件文件成功！\n")
		return
	}
	return
}

func token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, explorerToken)
}

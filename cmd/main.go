package main

import (
	"flag"
	"fmt"
	"git.iotserv.com/iotserv/client/config"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	h    bool
	id   string
	host string
	t    string
	p    string
)
var tkstr string = ""
var tkstr2 string = ""

func init() {
	flag.BoolVar(&h, "h", false, "show this help")
	flag.StringVar(&id, "id", uuid.Must(uuid.NewV4()).String(), "your client id")
	flag.StringVar(&host, "host", config.RegisterHost, "your server host")
	flag.StringVar(&t, "t", "", "set `token`")
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

func token(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, tkstr2)
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}
	configMode := config.ClientConfig{}
	_, err := os.Stat(config.ConfigFilePath)
	if err == nil && t == "" {
		//配置文件存在
		fmt.Println("找到了配置文件，开始使用上次的token连接")
		content, err := ioutil.ReadFile(config.ConfigFilePath)
		if err == nil {
			err = yaml.Unmarshal(content, &configMode)
			if err == nil {
				t = configMode.LastToken.ClientToken
			}
		}
	} else {
		//	配置文件不存在
		fmt.Println("配置文件不存在开始自动生成...")
	}
	//t = "eyJpZCI6InJ1bmlkdGVzdCIsImhzdCI6IjExOC44OS4xMDYuMjI2IiwicG90Ijo3MDAwLCJwZSI6MSwidGltIjoxNTIzNDU0NDUyLCJleHAiOjIwMDAwMDAsInNpZyI6IjBlNDYwZmZjNThiNTE5OTY1MDJhMDJkYTNkZTUwZjY3In0="
	if t != "" {
		_, err := services.RunNATManager(t)
		if err != nil {
			fmt.Printf(err.Error())
			fmt.Printf("登陆失败！请重新登陆。")
			os.Exit(0)
		} else {
			fmt.Printf("登陆成功！\n")
			fmt.Println("访问token：\n\n" + configMode.LastToken.ExplorerToken + "\n\n")
		}

	} else {
		//crypto.Salt = "abc"
		//b := make([]byte, 10)
		//rand.Read(b)
		//id := fmt.Sprintf("%x", b)
		if id == "" {
			id = uuid.Must(uuid.NewV4()).String()
		}
		fmt.Println("host:", host)
		tkstr, _ = crypto.GetToken(id, host, config.TcpPort, config.KcpPort, config.TlsPort, config.UdpApiPort, 1, 200000000000) //47.96.185.226，118.89.106.226
		fmt.Println(tkstr)
		_, err := services.RunNATManager(tkstr)
		if err != nil {
			fmt.Printf(err.Error())
			fmt.Printf("登陆失败！请重新登陆。")
			os.Exit(0)
		} else {
			fmt.Printf("登陆成功！\n")
		}
		//fmt.Printf(tkstr+"\n")
		tkstr2, _ = crypto.GetToken(id, host, config.TcpPort, config.KcpPort, config.TlsPort, config.UdpApiPort, 2, 200000000000) //118.89.106.226
		fmt.Printf("要想访问本内网，请用explorer使用以下token：\n" + tkstr2 + "\n")
		configMode.Common.Id = id
		configMode.Common.RegisterHost = host
		configMode.LastToken.ClientToken = tkstr
		configMode.LastToken.ExplorerToken = tkstr2
		configByte, err := yaml.Marshal(&configMode)
		if err != nil {
			fmt.Println(err.Error())
		}
		if ioutil.WriteFile(config.ConfigFilePath, configByte, 0644) == nil {
			fmt.Println("写入配置文件文件成功！\n")
		}
	}
	http.HandleFunc("/", token) //设置访问的路由
	//fmt.Printf("管理地址：http://127.0.0.1:"+p+"\n")
	err = http.ListenAndServe("127.0.0.1:"+p, nil) //设置监听的端口
	if err != nil {
		fmt.Printf("请检查端口" + p + "是否被占用")
	}
}

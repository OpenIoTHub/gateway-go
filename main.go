package main

import (
	"fmt"
	client "github.com/OpenIoTHub/gateway-go/v2/client"
	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/register"
	"github.com/OpenIoTHub/gateway-go/v2/services"
	"github.com/OpenIoTHub/gateway-go/v2/utils/login_utils"
	"github.com/OpenIoTHub/gateway-go/v2/utils/qr"
	"github.com/OpenIoTHub/gateway-go/v2/utils/str"
	utils_models "github.com/OpenIoTHub/utils/models"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var (
	version = ""
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	client.IsLibrary = false
	myApp := cli.NewApp()
	myApp.Name = "gateway-go"
	myApp.Usage = "-c [config file path]"
	myApp.Version = str.BuildVersion(version, commit, date, builtBy)
	myApp.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       config.ConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"GatewayConfigFilePath"},
					Destination: &config.ConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				config.InitConfigFile()
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test this command",
			Action: func(c *cli.Context) error {
				fmt.Println("ok")
				return nil
			},
		},
		//	TODO 与当前本机运行的服务进行通信查询一些信息，比如id，token，log等
		//	TODO 查询本机所在网络所包含的支持的服务
	}
	myApp.Flags = []cli.Flag{
		//TODO 应该设置工作目录，各组件共享
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       config.ConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"GatewayConfigFilePath"},
			Destination: &config.ConfigFilePath,
		},
		//token 登录
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"t"},
			Value:       config.GatewayLoginToken,
			Usage:       "login server by gateway token ",
			EnvVars:     []string{"GatewayLoginToken"},
			Destination: &config.GatewayLoginToken,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		str.PrintOpenIoTHubLogo()
		if config.GatewayLoginToken != "" {
			UseGateWayToken()
		} else {
			_, err := os.Stat(config.ConfigFilePath)
			if err != nil {
				config.InitConfigFile()
			}
			UseConfigFile()
		}
		go client.Run()
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
		return nil
	}
	register.RegisterService("localhost-gateway-go",
		"_http._tcp",
		"local",
		"localhost",
		client.HttpPort,
		[]string{"name=gateway-go", fmt.Sprintf("id=gateway-go@%s", uuid.Must(uuid.NewV4()).String()), "home-page=https://github.com/OpenIoTHub/gateway-go"},
		0,
		[]net.IP{net.ParseIP("127.0.0.1")},
		[]net.IP{},
	)
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

//export Run
func Run() {
	go client.Run()
}

func UseGateWayToken() {
	//使用服务器签发的Token登录
	err := services.GatewayManager.AddServer(config.GatewayLoginToken)
	if err != nil {
		log.Println(err.Error())
		log.Printf("登陆失败！请重新登陆。")
		return
	}
	log.Printf("登陆成功！\n")
}

func UseConfigFile() {
	//配置文件存在
	log.Println("使用的配置文件位置：", config.ConfigFilePath)
	content, err := os.ReadFile(config.ConfigFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &config.ConfigMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//找到了配置文件
	if len(config.ConfigMode.GatewayUUID) < 35 {
		config.ConfigMode.GatewayUUID = uuid.Must(uuid.NewV4()).String()
		err = config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if config.ConfigMode.LoginWithTokenMap == nil {
		config.ConfigMode.LoginWithTokenMap = make(map[string]string)
	}
	//解析日志配置
	writers := []io.Writer{}
	if config.ConfigMode.LogConfig.EnableStdout {
		writers = append(writers, os.Stdout)
	}
	if config.ConfigMode.LogConfig.LogFilePath != "" {
		f, err := os.OpenFile(config.ConfigMode.LogConfig.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writers = append(writers, f)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	//解析配置文件，解析服务器配置文件列表
	//解析登录token列表
	//如果CLI模式尚未登录自动登陆服务器并创建一个二维码
	if len(config.ConfigMode.LoginWithTokenMap) == 0 {
		err = login_utils.AutoLoginAndDisplayQRCode()
		if err != nil {
			log.Println(err)
		}
	}
	for _, v := range config.ConfigMode.LoginWithTokenMap {
		err = services.GatewayManager.AddServer(v)
		if err != nil {
			continue
		}
		// 通过gateway jwt(UUID)展示二维码
		tokenModel, err := utils_models.DecodeUnverifiedToken(v)
		if err != nil {
			return
		}
		err = qr.DisplayQRCodeById(tokenModel.RunId, tokenModel.Host)
		if err != nil {
			continue
		}
	}
}

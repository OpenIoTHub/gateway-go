package main

import (
	"fmt"
	client "github.com/OpenIoTHub/gateway-go/client"
	"github.com/OpenIoTHub/gateway-go/netservice/login"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = "time"
	builtBy = ""
)

func main() {
	login.Version = version
	myApp := cli.NewApp()
	myApp.Name = "gateway-go"
	myApp.Usage = "-c [config file path]"
	myApp.Version = buildVersion(version, commit, date, builtBy)
	myApp.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       services.ConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"GatewayConfigFilePath"},
					Destination: &services.ConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				services.InitConfigFile()
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
			Value:       services.ConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"GatewayConfigFilePath"},
			Destination: &services.ConfigFilePath,
		},
		//token 登录
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"t"},
			Value:       services.GatewayLoginToken,
			Usage:       "login server by gateway token ",
			EnvVars:     []string{"GatewayLoginToken"},
			Destination: &services.GatewayLoginToken,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		if services.GatewayLoginToken != "" {
			services.UseGateWayToken()
		} else {
			_, err := os.Stat(services.ConfigFilePath)
			if err != nil {
				services.InitConfigFile()
			}
			services.UseConfigFile()
		}
		go client.Run()
		for {
			time.Sleep(time.Hour)
		}
	}
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

//export Run
func Run() {
	go client.Run()
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}

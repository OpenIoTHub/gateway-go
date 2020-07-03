package main

import (
	"fmt"
	client "github.com/OpenIoTHub/gateway-go/client"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "gateway-go"
	myApp.Usage = "-c [config file path]"
	myApp.Version = fmt.Sprintf("%s(commit:%s,build on:%s,buildBy:%s)", client.Version, client.Commit, client.Date, client.BuiltBy)
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
			Value:       "",
			Usage:       "login server by gateway token ",
			EnvVars:     []string{"GatewayLoginToken"},
			Destination: &config.GatewayLoginToken,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		_, err := os.Stat(config.ConfigFilePath)
		if err != nil {
			config.InitConfigFile(client.ConfigMode)
		}
		config.UseConfigFile(client.ConfigMode)
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

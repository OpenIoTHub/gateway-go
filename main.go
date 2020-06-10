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
			Name:    "config",
			Aliases: []string{"c"},
			Value:   config.Setting["configFilePath"],
			Usage:   "config file path",
			EnvVars: []string{"GatewayConfigFilePath"},
		},
	}
	myApp.Action = func(c *cli.Context) error {
		if c.String("config") != "" {
			config.Setting["configFilePath"] = c.String("config")
		}
		_, err := os.Stat(config.Setting["configFilePath"])
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

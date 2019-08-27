package main

import (
	"fmt"
	client "git.iotserv.com/iotserv/client/android"
	"git.iotserv.com/iotserv/client/config"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/urfave/cli"
	"os"
	"time"
)

func init() {
	go client.Run()
}

func main() {
	configMode := models.ClientConfig{}

	myApp := cli.NewApp()
	myApp.Name = "client"
	myApp.Usage = "-c [配置文件路径]"
	myApp.Version = "v1.0.1"
	myApp.Flags = []cli.Flag{
		//TODO 应该设置工作目录，各组件共享
		cli.StringFlag{
			Name:   "config,c",
			Value:  config.Setting["configFilePath"],
			Usage:  "配置文件路径",
			EnvVar: "Config_File_Path",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		if c.String("config") != "" {
			config.Setting["configFilePath"] = c.String("config")
		}
		_, err := os.Stat(config.Setting["configFilePath"])
		if err != nil {
			config.InitConfigFile(configMode)
		} else {
			config.UseConfigFile(configMode)
		}
		for {
			time.Sleep(time.Hour)
		}
	}
	err := myApp.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}

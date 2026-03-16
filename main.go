package main

import (
	"fmt"
	"log"
	"os"

	client "github.com/OpenIoTHub/gateway-go/v2/client"
	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/info"
	"github.com/OpenIoTHub/gateway-go/v2/services"
	"github.com/urfave/cli/v2"
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
	myApp.Version = info.BuildVersion(version, commit, date, builtBy)
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
	}
	myApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       config.ConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"GatewayConfigFilePath"},
			Destination: &config.ConfigFilePath,
		},
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"t"},
			Value:       config.GatewayLoginToken,
			Usage:       "login server by gateway token",
			EnvVars:     []string{"GatewayLoginToken"},
			Destination: &config.GatewayLoginToken,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		info.PrintLogo()
		if config.GatewayLoginToken != "" {
			services.StartWithToken(config.GatewayLoginToken)
		} else {
			if _, err := os.Stat(config.ConfigFilePath); err != nil {
				config.InitConfigFile()
			}
			services.StartWithConfigFile()
		}
		go client.Run()
		select {}
	}
	if err := myApp.Run(os.Args); err != nil {
		log.Println(err.Error())
	}
}

//export Run
func Run() {
	go client.Run()
}

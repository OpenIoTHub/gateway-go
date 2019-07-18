package main

import (
	"fmt"
	"git.iotserv.com/iotserv/client/config"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	configMode := models.ClientConfig{}

	myApp := cli.NewApp()
	myApp.Name = "client"
	myApp.Usage = "-c [配置文件路径]"
	myApp.Version = "v1.0.1"
	myApp.Flags = []cli.Flag{
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
			initConfigFile(configMode)
			return err
		}
		useConfigFile(configMode)
		for {
			time.Sleep(time.Hour)
		}
	}
	err := myApp.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
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

func initConfigFile(configMode models.ClientConfig) {
	fmt.Println("没有找到配置文件：", config.Setting["configFilePath"])
	fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
	//	生成配置文件模板
	port, _ := strconv.Atoi(config.Setting["apiPort"])
	configMode.ExplorerTokenHttpPort = port
	configMode.Server.ConnectionType = "tcp"
	configMode.Server.ServerHost = "guonei.nat-cloud.com"
	configMode.Server.TcpPort = 34320
	configMode.Server.KcpPort = 34320
	configMode.Server.UdpApiPort = 34321
	configMode.Server.TlsPort = 34321
	configMode.Server.LoginKey = "HLLdsa544&*S"
	//configMode.Server.ServerHost = "netipcam.com"
	//configMode.Server.TcpPort = 5555
	//configMode.Server.KcpPort = 5555
	//configMode.Server.UdpApiPort = 6666
	//configMode.Server.TlsPort = 6666
	//configMode.Server.LoginKey = "kasan@KASAN5555"
	configMode.LastId = uuid.Must(uuid.NewV4()).String()
	err := writeConfigFile(configMode, config.Setting["configFilePath"])
	if err == nil {
		fmt.Println("由于没有找到配置文件，已经为你生成配置文件（模板），位置：", config.Setting["configFilePath"])
		fmt.Println("你可以手动修改上述配置文件后再运行！")
		return
	}
	fmt.Println("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
	fmt.Println(err.Error())
}

func useConfigFile(configMode models.ClientConfig) {
	//配置文件存在
	fmt.Println("使用的配置文件位置：", config.Setting["configFilePath"])
	content, err := ioutil.ReadFile(config.Setting["configFilePath"])
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
	config.Setting["clientToken"], err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 1, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = services.RunNATManager(configMode.Server.LoginKey, config.Setting["clientToken"])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		os.Exit(0)
	}
	fmt.Printf("登陆成功！\n")
	config.Setting["explorerToken"], err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 2, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("访问token：\n\n" + config.Setting["explorerToken"] + "\n\n")
	err = writeConfigFile(configMode, config.Setting["configFilePath"])
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("你也可以访问：http://127.0.0.1:%s/查看访问token\n", config.Setting["apiPort"])
}

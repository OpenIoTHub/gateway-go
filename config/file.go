package config

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/utils/crypto"
	"github.com/OpenIoTHub/utils/models"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

//将配置写入指定的路径的文件
func WriteConfigFile(configMode models.ClientConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		fmt.Println("写入配置文件文件成功！")
		return
	}
	return
}

func InitConfigFile(configMode models.ClientConfig) {
	fmt.Println("没有找到配置文件：", Setting["configFilePath"])
	fmt.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
	//	生成配置文件模板
	port, _ := strconv.Atoi(Setting["apiPort"])
	configMode.ExplorerTokenHttpPort = port
	configMode.Server.ConnectionType = "tcp"

	configMode.Server.ServerHost = "guonei.nat-cloud.com"
	configMode.Server.TcpPort = 34320
	configMode.Server.KcpPort = 34320
	configMode.Server.UdpApiPort = 34321
	configMode.Server.TlsPort = 34321
	configMode.Server.LoginKey = "HLLdsa544&*S"

	configMode.LastId = uuid.Must(uuid.NewV4()).String()
	err := WriteConfigFile(configMode, Setting["configFilePath"])
	if err == nil {
		fmt.Println("由于没有找到配置文件，已经为你生成配置文件（模板），位置：", Setting["configFilePath"])
		fmt.Println("你可以手动修改上述配置文件后再运行！")
		return
	}
	fmt.Println("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
	fmt.Println(err.Error())
}

func UseConfigFile(configMode models.ClientConfig) {
	//配置文件存在
	fmt.Println("使用的配置文件位置：", Setting["configFilePath"])
	content, err := ioutil.ReadFile(Setting["configFilePath"])
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
	Setting["clientToken"], err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 1, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = services.RunNATManager(configMode.Server.LoginKey, Setting["clientToken"])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	}
	fmt.Printf("登陆成功！\n")
	Setting["explorerToken"], err = crypto.GetToken(configMode.Server.LoginKey, configMode.LastId, configMode.Server.ServerHost, configMode.Server.TcpPort,
		configMode.Server.KcpPort, configMode.Server.TlsPort, configMode.Server.UdpApiPort, 2, 200000000000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("访问token：\n\n" + Setting["explorerToken"] + "\n\n")
	err = WriteConfigFile(configMode, Setting["configFilePath"])
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("你也可以访问：http://127.0.0.1:%s/查看访问token\n", Setting["apiPort"])
	Loged = true
}

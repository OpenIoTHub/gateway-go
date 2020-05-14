package config

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/utils/models"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

//将配置写入指定的路径的文件
func WriteConfigFile(configMode *models.GatewayConfig, path string) (err error) {
	configByte, err := yaml.Marshal(configMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		log.Println("写入配置文件文件成功！")
		return
	}
	return
}

func InitConfigFile(configMode *models.GatewayConfig) {
	log.Println("没有找到配置文件：", Setting["configFilePath"])
	log.Println("开始生成默认的空白配置文件，请填写配置文件后重复运行本程序")
	//	生成配置文件模板
	port, _ := strconv.Atoi(Setting["gRpcPort"])
	configMode.GrpcPort = port
	configMode.ConnectionType = "tcp"

	configMode.Server.ServerHost = "guonei.nat-cloud.com"
	configMode.Server.TcpPort = 34320
	configMode.Server.KcpPort = 34320
	configMode.Server.UdpApiPort = 34321
	configMode.Server.KcpApiPort = 34322
	configMode.Server.TlsPort = 34321
	configMode.Server.GrpcPort = 34322
	configMode.Server.LoginKey = "HLLdsa544&*S"

	configMode.LastId = uuid.Must(uuid.NewV4()).String()
	err := os.MkdirAll(filepath.Dir(Setting["configFilePath"]), 0644)
	if err != nil {
		return
	}
	err = WriteConfigFile(configMode, Setting["configFilePath"])
	if err == nil {
		log.Println("由于没有找到配置文件，已经为你生成配置文件（模板），位置：", Setting["configFilePath"])
		log.Println("你可以手动修改上述配置文件后再运行！")
		return
	}
	log.Println("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
	log.Println(err.Error())
}

func UseConfigFile(configMode *models.GatewayConfig) {
	//配置文件存在
	log.Println("使用的配置文件位置：", Setting["configFilePath"])
	content, err := ioutil.ReadFile(Setting["configFilePath"])
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &configMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//找到了配置文件
	if len(configMode.LastId) < 35 {
		configMode.LastId = uuid.Must(uuid.NewV4()).String()
	}
	Setting["GateWayToken"], err = models.GetToken(configMode, 1, 200000000000)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = services.RunNATManager(configMode.Server.LoginKey, Setting["GateWayToken"])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	}
	fmt.Printf("登陆成功！\n")
	Setting["OpenIoTHubToken"], err = models.GetToken(configMode, 2, 200000000000)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("访问token：\n\n" + Setting["OpenIoTHubToken"] + "\n\n")
	err = WriteConfigFile(configMode, Setting["configFilePath"])
	if err != nil {
		log.Println(err.Error())
	}
	Loged = true
}

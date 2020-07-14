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
)

//将配置写入指定的路径的文件
func WriteConfigFile(ConfigMode *models.GatewayConfig, path string) (err error) {
	configByte, err := yaml.Marshal(ConfigMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if ioutil.WriteFile(path, configByte, 0644) == nil {
		return
	}
	return
}

func InitConfigFile() {
	//	生成配置文件模板
	err := os.MkdirAll(filepath.Dir(ConfigFilePath), 0644)
	if err != nil {
		return
	}
	err = WriteConfigFile(ConfigMode, ConfigFilePath)
	if err == nil {
		fmt.Println("config created")
		return
	}
	log.Println("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
	log.Println(err.Error())
}

func UseConfigFile() {
	//配置文件存在
	log.Println("使用的配置文件位置：", ConfigFilePath)
	content, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(content, &ConfigMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//找到了配置文件
	if len(ConfigMode.LastId) < 35 {
		ConfigMode.LastId = uuid.Must(uuid.NewV4()).String()
	}
	Setting["GateWayToken"], err = models.GetToken(ConfigMode, 1, 200000000000)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = services.RunNATManager(ConfigMode.Server.LoginKey, Setting["GateWayToken"])
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	}
	fmt.Printf("登陆成功！\n")
	Setting["OpenIoTHubToken"], err = models.GetToken(ConfigMode, 2, 200000000000)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("访问token：\n\n" + Setting["OpenIoTHubToken"] + "\n\n")
	err = WriteConfigFile(ConfigMode, ConfigFilePath)
	if err != nil {
		log.Println(err.Error())
	}
	Loged = true
}

func UseGateWayToken() {
	//使用服务器签发的Token登录
	err := services.RunNATManager(ConfigMode.Server.LoginKey, GatewayLoginToken)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Printf("登陆失败！请重新登陆。")
		return
	}
	fmt.Printf("登陆成功！\n")
	Loged = true
}

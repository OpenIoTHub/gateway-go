package config

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/models"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

// 将配置写入指定的路径的文件
func WriteConfigFile(ConfigMode *models.GatewayConfig, path string) (err error) {
	configByte, err := yaml.Marshal(ConfigMode)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if os.WriteFile(path, configByte, 0644) == nil {
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

//func UseGateWayToken() {
//	//使用服务器签发的Token登录
//	err := GatewayManager.AddServer(GatewayLoginToken)
//	if err != nil {
//		log.Printf(err.Error())
//		log.Printf("登陆失败！请重新登陆。")
//		return
//	}
//	log.Printf("登陆成功！\n")
//}

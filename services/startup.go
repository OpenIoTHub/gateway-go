package services

import (
	"log"
	"os"

	"github.com/OpenIoTHub/gateway-go/v2/config"
	"github.com/OpenIoTHub/gateway-go/v2/utils/qr"
	utils_models "github.com/OpenIoTHub/utils/v2/models"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
)

func StartWithToken(token string) {
	if err := GatewayManager.AddServer(token); err != nil {
		log.Printf("登录失败: %v，请重新登录", err)
		return
	}
	log.Println("登录成功！")
}

func StartWithConfigFile() {
	log.Println("使用的配置文件位置：", config.ConfigFilePath)
	content, err := os.ReadFile(config.ConfigFilePath)
	if err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return
	}
	if err := yaml.Unmarshal(content, &config.ConfigMode); err != nil {
		log.Printf("解析配置文件失败: %v", err)
		return
	}
	if len(config.ConfigMode.GatewayUUID) < 35 {
		config.ConfigMode.GatewayUUID = uuid.Must(uuid.NewV4()).String()
		if err := config.WriteConfigFile(config.ConfigMode, config.ConfigFilePath); err != nil {
			log.Printf("写入配置文件失败: %v", err)
		}
	}
	if config.ConfigMode.LoginWithTokenMap == nil {
		config.ConfigMode.LoginWithTokenMap = make(map[string]string)
	}
	config.SetupLogging(config.ConfigMode.LogConfig)

	if len(config.ConfigMode.LoginWithTokenMap) == 0 {
		if err := AutoLoginAndDisplayQRCode(); err != nil {
			log.Printf("自动登录失败: %v", err)
		}
	}
	for _, token := range config.ConfigMode.LoginWithTokenMap {
		if err := GatewayManager.AddServer(token); err != nil {
			log.Printf("添加服务器失败: %v", err)
			continue
		}
		tokenModel, err := utils_models.DecodeUnverifiedToken(token)
		if err != nil {
			log.Printf("解析token失败: %v", err)
			continue
		}
		if err := qr.DisplayQRCodeById(tokenModel.RunId, tokenModel.Host); err != nil {
			log.Printf("显示二维码失败: %v", err)
			continue
		}
	}
}

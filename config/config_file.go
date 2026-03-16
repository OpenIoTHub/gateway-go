package config

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/v2/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

// 将配置写入指定的路径的文件
func WriteConfigFile(ConfigMode *models.GatewayConfig, path string) error {
	configByte, err := yaml.Marshal(ConfigMode)
	if err != nil {
		return fmt.Errorf("marshal config failed: %w", err)
	}
	if err := os.WriteFile(path, configByte, 0644); err != nil {
		return fmt.Errorf("write config file failed: %w", err)
	}
	return nil
}

func InitConfigFile() {
	//	生成配置文件模板
	dir := filepath.Dir(ConfigFilePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建配置目录失败: %v", err)
			return
		}
	}
	if err := WriteConfigFile(ConfigMode, ConfigFilePath); err != nil {
		log.Println("写入配置文件模板出错，请检查本程序是否具有写入权限！或者手动创建配置文件。")
		log.Println(err.Error())
		return
	}
	fmt.Println("config created")
}

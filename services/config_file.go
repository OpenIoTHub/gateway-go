package services

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/models"
	utils_models "github.com/OpenIoTHub/utils/models"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io"
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

func UseConfigFile() {
	//配置文件存在
	log.Println("使用的配置文件位置：", ConfigFilePath)
	content, err := os.ReadFile(ConfigFilePath)
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
	if len(ConfigMode.GatewayUUID) < 35 {
		ConfigMode.GatewayUUID = uuid.Must(uuid.NewV4()).String()
		err = WriteConfigFile(ConfigMode, ConfigFilePath)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if ConfigMode.LoginWithTokenMap == nil {
		ConfigMode.LoginWithTokenMap = make(map[string]string)
	}
	//解析日志配置
	writers := []io.Writer{}
	if ConfigMode.LogConfig.EnableStdout {
		writers = append(writers, os.Stdout)
	}
	if ConfigMode.LogConfig.LogFilePath != "" {
		f, err := os.OpenFile(ConfigMode.LogConfig.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		writers = append(writers, f)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	//解析配置文件，解析服务器配置文件列表
	//解析登录token列表
	//如果CLI模式尚未登录自动登陆服务器并创建一个二维码
	if len(ConfigMode.LoginWithTokenMap) == 0 {
		err = autoLoginAndDisplayQRCode()
		if err != nil {
			log.Println(err)
		}
	}
	for _, v := range ConfigMode.LoginWithTokenMap {
		err = GatewayManager.AddServer(v)
		if err != nil {
			continue
		}
		// 通过gateway jwt(UUID)展示二维码
		tokenModel, err := utils_models.DecodeUnverifiedToken(v)
		if err != nil {
			log.Printf("DecodeUnverifiedToken:", err.Error())
			return
		}
		err = displayQRCodeById(tokenModel.RunId)
		if err != nil {
			continue
		}
	}
}

func UseGateWayToken() {
	//使用服务器签发的Token登录
	err := GatewayManager.AddServer(GatewayLoginToken)
	if err != nil {
		log.Printf(err.Error())
		log.Printf("登陆失败！请重新登陆。")
		return
	}
	log.Printf("登陆成功！\n")
}

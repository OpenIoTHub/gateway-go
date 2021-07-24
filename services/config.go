package services

import (
	"fmt"
	"github.com/OpenIoTHub/gateway-go/models"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
)

const ConfigFileName = "gateway-go.yaml"

var ConfigFilePath = fmt.Sprintf("%s%s", "./", ConfigFileName)

var GatewayLoginToken = ""

const GRpcAddr = ""
const GrpcPort = 0

var ConfigMode = &models.GatewayConfig{
	GatewayUUID: uuid.Must(uuid.NewV4()).String(),
	LogConfig: &models.LogConfig{
		EnableStdout: true,
		LogFilePath:  "",
	},
	LoginWithTokenMap: map[string]string{},
}

func init() {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
}

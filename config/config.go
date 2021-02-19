package config

import (
	"fmt"
	"github.com/OpenIoTHub/utils/models"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
)

//var Setting = make(map[string]string)

var Loged = false

const ConfigFileName = "gateway-go.yaml"

var ConfigFilePath = fmt.Sprintf("%s%s", "./", ConfigFileName)

var GatewayLoginToken = ""
var OpenIoTHubToken = ""

const GRpcAddr = "0.0.0.0"
const GrpcPort = 0

var ConfigMode = &models.GatewayConfig{
	GatewayUUID: uuid.Must(uuid.NewV4()).String(),
	LogConfig: &models.LogConfig{
		EnableStdout: true,
		LogFilePath:  "",
	},
	LoginWithTokenList: []string{},
	LoginWithServerConf: []*models.LoginWithServer{
		{
			LastId:         uuid.Must(uuid.NewV4()).String(),
			ConnectionType: "tcp",
			Server: &models.Srever{
				ServerHost: "guonei.nat-cloud.com",
				TcpPort:    34320,
				KcpPort:    34320,
				UdpApiPort: 34321,
				KcpApiPort: 34322,
				TlsPort:    34321,
				GrpcPort:   34322,
				LoginKey:   "HLLdsa544&*S",
			},
		},
	},
}

func init() {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
}

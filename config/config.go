package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var Setting = make(map[string]string)

var Loged = false

var ConfigFileName = "gateway.yaml"
var ConfigFilePath = fmt.Sprintf("%s%s", "./", ConfigFileName)
var GatewayLoginToken = ""

func init() {
	//是否是snapcraft应用，如果是则从snapcraft指定的工作目录保存配置文件
	appDataPath, havaAppDataPath := os.LookupEnv("SNAP_USER_DATA")
	if havaAppDataPath {
		ConfigFilePath = filepath.Join(appDataPath, ConfigFileName)
	}
	Setting["gRpcAddr"] = "0.0.0.0"
	Setting["gRpcPort"] = "1082"
	Setting["GateWayToken"] = ""
	Setting["OpenIoTHubToken"] = ""
	if runtime.GOOS == "android" {

	} else {

	}
}

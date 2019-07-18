package config

import (
	"runtime"
)

var Setting = make(map[string]string)

func init() {
	Setting["configFilePath"] = "./client.yaml"
	Setting["apiPort"] = "1082"
	Setting["clientToken"] = ""
	Setting["explorerToken"] = ""
	if runtime.GOOS == "android" {

	} else {

	}
}

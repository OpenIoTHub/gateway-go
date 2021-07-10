package models

//网关的配置文件
type GatewayConfig struct {
	GatewayUUID        string     `json:"uuid"`
	LogConfig          *LogConfig `json:"log"`
	LoginWithTokenList []string   `json:"tokens"`
}

type LogConfig struct {
	EnableStdout bool
	LogFilePath  string
}

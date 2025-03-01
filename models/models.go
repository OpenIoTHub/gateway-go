package models

// 网关的配置文件
type GatewayConfig struct {
	GatewayUUID       string            `json:"uuid"`
	LogConfig         *LogConfig        `json:"log"`
	HttpServicePort   int               `json:"http_service_port" yaml:"http_service_port"`
	LoginWithTokenMap map[string]string `json:"tokens"`
}

type LogConfig struct {
	EnableStdout bool
	LogFilePath  string
}

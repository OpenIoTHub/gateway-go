package config

import (
	"io"
	"log"
	"os"

	"github.com/OpenIoTHub/gateway-go/v2/models"
)

// SetupLogging 根据配置初始化日志输出目标
func SetupLogging(logConfig *models.LogConfig) {
	if logConfig == nil {
		return
	}
	var writers []io.Writer
	if logConfig.EnableStdout {
		writers = append(writers, os.Stdout)
	}
	if logConfig.LogFilePath != "" {
		f, err := os.OpenFile(logConfig.LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		writers = append(writers, f)
	}
	if len(writers) > 0 {
		log.SetOutput(io.MultiWriter(writers...))
	}
}

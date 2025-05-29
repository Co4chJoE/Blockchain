// utils/logger.go
package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	// 设置日志级别为 Info 及以上
	Logger.SetLevel(logrus.InfoLevel)

	// 输出到控制台和文件
	file, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
	} else {
		Logger.Warn("无法打开日志文件，将只输出到控制台")
	}
}

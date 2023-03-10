package logger

import (
	"github.com/antnzr/chat-go/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func getLogConfig() zap.Config {
	config, _ := config.LoadConfig(".")
	env := config.GinMode
	if env == gin.ReleaseMode {
		config := zap.NewProductionConfig()
		zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
		return config
	}
	return zap.NewDevelopmentConfig()
}

func init() {
	var err error
	config := getLogConfig()
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	defer zapLog.Sync()
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}

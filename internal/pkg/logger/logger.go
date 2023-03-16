package logger

import (
	"github.com/antnzr/chat-go/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func init() {
	mainConfig, _ := config.LoadConfig(".")

	var err error
	config := getLogConfig(mainConfig)
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return zapLog
}

func Flush() error {
	return zapLog.Sync()
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

func Fatality(fields ...zap.Field) {
	zapLog.Fatal("", fields...)
}

func getLogConfig(config config.Config) zap.Config {
	var zapConfig zap.Config
	if config.GinMode == gin.ReleaseMode {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.EncoderConfig.TimeKey = "@ts"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapConfig
}

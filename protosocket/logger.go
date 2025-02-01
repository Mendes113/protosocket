package protosocket

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// GetLogger retorna a instância do logger
func GetLogger() *zap.Logger {
	return logger
}

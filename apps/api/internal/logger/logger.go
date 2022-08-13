package logger

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/env"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildBaseLogger(cfg *config.Config) *zap.Logger {
	var zc zap.Config
	if cfg.Env == env.EnvDevelopment {
		zc = zapdriver.NewDevelopmentConfig()
		zc.Encoding = "json"
	} else {
		zc = zapdriver.NewProductionConfig()
		zc.EncoderConfig.TimeKey = zapcore.OmitKey
	}
	logger, _ := zc.Build()
	return logger
}

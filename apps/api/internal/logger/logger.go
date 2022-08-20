package logger

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/go/env"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildBaseLogger(cfg *config.Config) *zap.Logger {
	var zc zap.Config
	if cfg.Env == env.EnvDevelopment {
		zc = zapdriver.NewDevelopmentConfig()
	} else if cfg.Env == env.EnvStaging {
		zc = zapdriver.NewProductionConfig()
		zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zc.EncoderConfig.TimeKey = zapcore.OmitKey
	} else {
		zc = zapdriver.NewProductionConfig()
		zc.EncoderConfig.TimeKey = zapcore.OmitKey
	}
	logger, _ := zc.Build()
	return logger
}

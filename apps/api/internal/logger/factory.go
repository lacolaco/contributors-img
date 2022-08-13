package logger

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/apiclient"
	"contrib.rocks/libs/goutils/env"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
)

type LoggerFactory interface {
	Logger(name string) Logger
	Close()
}

func NewLoggerFactory(cfg *config.Config) LoggerFactory {
	if cfg.ProjectID() != "" {
		return newCloudLoggingLoggerFactory(apiclient.NewLoggingClient(cfg.ProjectID()))
	} else {
		return newStdLoggerFactory()
	}
}

func buildBaseLogger(cfg *config.Config) *zap.Logger {
	var zc zap.Config
	if cfg.Env == env.EnvDevelopment {
		zc = zapdriver.NewDevelopmentConfig()
		zc.Encoding = "json"
	} else {
		zc = zapdriver.NewProductionConfig()
	}
	logger, _ := zc.Build()
	return logger
}

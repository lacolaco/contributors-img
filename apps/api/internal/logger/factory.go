package logger

import (
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/apiclient"
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

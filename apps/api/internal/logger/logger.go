package logger

import (
	"context"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/apiclient"
)

type Logger interface {
	ContextWithLogger(parent context.Context) context.Context
	Log(ctx context.Context, e logging.Entry)
	Debug(ctx context.Context, e logging.Entry)
	Info(ctx context.Context, e logging.Entry)
	Error(ctx context.Context, e logging.Entry)
}

var loggerFactory *factory

func NewLogger(name string) Logger {
	return loggerFactory.Logger(name)
}

func InitLoggerFactory(cfg *config.Config) func() {
	cleanup := func() {}
	if cfg.ProjectID() != "" {
		loggingClient := apiclient.NewLoggingClient(cfg.ProjectID())
		cleanup = func() { loggingClient.Close() }
		loggerFactory = newLoggerFactory(cfg, loggingClient)
	} else {
		loggerFactory = newLoggerFactory(cfg, nil)
	}

	return cleanup
}

func NewEntry(o any) logging.Entry {
	return logging.Entry{
		Payload: o,
	}
}

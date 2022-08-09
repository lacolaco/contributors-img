package logger

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/env"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const (
	loggerContextKey contextKey = contextKey("logger")
)

var _ Logger = &loggerWrapper{}

type loggerWrapper struct {
	cfg    *config.Config
	logger *logging.Logger
}

func (l *loggerWrapper) Log(c context.Context, entry logging.Entry) {
	if entry.Trace == "" {
		entry.Trace = l.trace(c)
	}
	if l.cfg.Env == env.EnvDevelopment {
		log.Printf("[%s] %+v", entry.Severity, entry.Payload)
	}
	l.logger.Log(entry)
}

func (l *loggerWrapper) Debug(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Debug
	l.Log(c, entry)
}

func (l *loggerWrapper) Info(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Info
	l.Log(c, entry)
}

func (l *loggerWrapper) Error(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Error
	l.Log(c, entry)
}

func (l *loggerWrapper) trace(c context.Context) string {
	traceID := trace.SpanContextFromContext(c).TraceID()
	if traceID.IsValid() {
		return fmt.Sprintf("/projects/%s/traces/%s", l.cfg.ProjectID(), traceID.String())
	}
	return ""
}

func (l *loggerWrapper) ContextWithLogger(c context.Context) context.Context {
	return context.WithValue(c, loggerContextKey, l)
}

func fromContext(c context.Context) Logger {
	logger := c.Value(loggerContextKey).(Logger)
	return logger
}

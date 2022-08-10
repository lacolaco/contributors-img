package logger

import (
	"context"
	"encoding/json"
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

var _ Logger = &cloudLoggingLogger{}

type cloudLoggingLogger struct {
	cfg    *config.Config
	logger *logging.Logger
}

func (l *cloudLoggingLogger) Log(c context.Context, entry logging.Entry) {
	if entry.Trace == "" {
		entry.Trace = l.trace(c)
	}
	if l.cfg.Env == env.EnvDevelopment {
		bytes, err := json.Marshal(entry.Payload)
		if err != nil {
			log.Printf("Error marshalling entry: %v", err)
			return
		}
		log.Printf("[%s] %s", entry.Severity, string(bytes))
	}
	l.logger.Log(entry)
}

func (l *cloudLoggingLogger) Debug(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Debug
	l.Log(c, entry)
}

func (l *cloudLoggingLogger) Info(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Info
	l.Log(c, entry)
}

func (l *cloudLoggingLogger) Error(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Error
	l.Log(c, entry)
}

func (l *cloudLoggingLogger) trace(c context.Context) string {
	traceID := trace.SpanContextFromContext(c).TraceID()
	if traceID.IsValid() {
		return fmt.Sprintf("/projects/%s/traces/%s", l.cfg.ProjectID(), traceID.String())
	}
	return ""
}

func (l *cloudLoggingLogger) ContextWithLogger(c context.Context) context.Context {
	return context.WithValue(c, loggerContextKey, l)
}

func fromContext(c context.Context) Logger {
	logger := c.Value(loggerContextKey).(Logger)
	return logger
}

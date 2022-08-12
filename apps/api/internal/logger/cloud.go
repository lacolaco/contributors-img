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

type cloudLoggingLoggerFactory struct {
	client *logging.Client

	loggerCache map[string]*cloudLoggingLogger
}

func newCloudLoggingLoggerFactory(client *logging.Client) LoggerFactory {
	return &cloudLoggingLoggerFactory{client, make(map[string]*cloudLoggingLogger)}
}

func (l *cloudLoggingLoggerFactory) Logger(name string) Logger {
	if l.loggerCache[name] == nil {
		l.loggerCache[name] = &cloudLoggingLogger{l.client.Logger(name)}
	}
	return l.loggerCache[name]
}

func (l *cloudLoggingLoggerFactory) Close() {
	l.client.Close()
}

type cloudLoggingLogger struct {
	logger *logging.Logger
}

func (l *cloudLoggingLogger) Log(c context.Context, entry logging.Entry) {
	l.log(c, entry)
}

func (l *cloudLoggingLogger) Debug(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Debug
	l.log(c, entry)
}

func (l *cloudLoggingLogger) Info(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Info
	l.log(c, entry)
}

func (l *cloudLoggingLogger) Error(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Error
	l.log(c, entry)
}

func (l *cloudLoggingLogger) Close() {
	l.logger.Flush()
}

func (l *cloudLoggingLogger) log(c context.Context, entry logging.Entry) {
	cfg := config.FromContext(c)
	if entry.Labels == nil {
		entry.Labels = make(map[string]string)
	}
	entry.Labels["environment"] = string(cfg.Env)

	sc := trace.SpanContextFromContext(c)
	if sc.IsValid() {
		entry.Trace = fmt.Sprintf("projects/%s/traces/%s", cfg.ProjectID(), sc.TraceID().String())
		entry.SpanID = sc.SpanID().String()
		entry.TraceSampled = sc.IsSampled()
	}

	if cfg.Env == env.EnvDevelopment {
		bytes, err := json.Marshal(entry.Payload)
		if err != nil {
			log.Printf("Error marshalling entry: %v", err)
			return
		}
		log.Printf("[%s] %s", entry.Severity, string(bytes))
	}
	l.logger.Log(entry)
}

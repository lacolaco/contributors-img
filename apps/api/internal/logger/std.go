package logger

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"cloud.google.com/go/logging"
)

var _ Logger = &stdLogger{}

type stdLogger struct {
	name string
}

func newStdLogger(name string) *stdLogger {
	return &stdLogger{name}
}

func (l *stdLogger) Log(c context.Context, entry logging.Entry) {
	l.printStandardLog(entry)
}
func (l *stdLogger) Debug(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Debug
	l.printStandardLog(entry)
}
func (l *stdLogger) Info(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Info
	l.printStandardLog(entry)
}
func (l *stdLogger) Error(c context.Context, entry logging.Entry) {
	entry.Severity = logging.Error
	l.printStandardLog(entry)
}
func (l *stdLogger) ContextWithLogger(c context.Context) context.Context {
	return context.WithValue(c, loggerContextKey, l)
}

func (l *stdLogger) printStandardLog(entry logging.Entry) {
	bytes, err := json.Marshal(entry.Payload)
	if err != nil {
		log.Printf("Error marshalling entry: %v", err)
		return
	}
	log.Printf("[%s][%s] %s", l.name, strings.ToUpper(entry.Severity.String()), string(bytes))
}

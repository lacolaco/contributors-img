package logger

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/logging"
)

var _ Logger = &stdLogger{}

type stdLogger struct{}

func (l *stdLogger) Log(c context.Context, entry logging.Entry) {
	printStandardLog(entry)
}
func (l *stdLogger) Debug(c context.Context, entry logging.Entry) {
	printStandardLog(entry)
}
func (l *stdLogger) Info(c context.Context, entry logging.Entry) {
	printStandardLog(entry)
}
func (l *stdLogger) Error(c context.Context, entry logging.Entry) {
	printStandardLog(entry)
}
func (l *stdLogger) ContextWithLogger(c context.Context) context.Context {
	return context.WithValue(c, loggerContextKey, l)
}

func printStandardLog(entry logging.Entry) {
	bytes, err := json.Marshal(entry.Payload)
	if err != nil {
		log.Printf("Error marshalling entry: %v", err)
		return
	}
	log.Printf("[%s] %s", entry.Severity, string(bytes))
}

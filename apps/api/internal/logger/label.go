package logger

import (
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap/zapcore"
)

const (
	logGroupKey = "app.group"
)

func LogGroupLabel(group string) zapcore.Field {
	return zapdriver.Labels(zapdriver.Label(logGroupKey, group))
}

func Label(key, value string) zapcore.Field {
	return zapdriver.Labels(zapdriver.Label(key, value))
}

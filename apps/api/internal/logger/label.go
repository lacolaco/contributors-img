package logger

import (
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap/zapcore"
)

const (
	logGroupIDKey = "groupId"
)

func LogGroup(groupID string) zapcore.Field {
	return zapdriver.Labels(zapdriver.Label(logGroupIDKey, groupID))
}

func Label(key, value string) zapcore.Field {
	return zapdriver.Labels(zapdriver.Label(key, value))
}

package logger

import (
	"context"

	"cloud.google.com/go/logging"
)

type Logger interface {
	Log(ctx context.Context, e logging.Entry)
	Debug(ctx context.Context, e logging.Entry)
	Info(ctx context.Context, e logging.Entry)
	Error(ctx context.Context, e logging.Entry)
}

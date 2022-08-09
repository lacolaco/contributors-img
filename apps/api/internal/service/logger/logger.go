package logger

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/apps/api/internal/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	ContextWithLogger(parent context.Context) context.Context
	Log(ctx context.Context, e logging.Entry)
	Debug(ctx context.Context, e logging.Entry)
	Info(ctx context.Context, e logging.Entry)
	Error(ctx context.Context, e logging.Entry)
}

type LoggerFactory interface {
	Logger(name string) Logger
}

func NewLoggerFactory(cfg *config.Config, l *logging.Client) LoggerFactory {
	return newLoggerFactory(cfg, l)
}

func NewEntry(o any) logging.Entry {
	return logging.Entry{
		Payload: o,
	}
}

func Middleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.ContextWithLogger(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if !traceID.IsValid() {
			ctx, span := tracing.DefaultTracer.Start(ctx, "api.http")
			defer span.End()
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
		logger.Log(c.Request.Context(), logging.Entry{
			HTTPRequest: &logging.HTTPRequest{
				Request: c.Request,
			},
			Timestamp: time.Now(),
			Severity:  logging.Info,
			Payload: map[string]any{
				"status": fmt.Sprintf("%d", c.Writer.Status()),
				"method": c.Request.Method,
				"url":    c.Request.URL.String(),
			},
		})
	}
}

func FromContext(c context.Context) Logger {
	return fromContext(c)
}

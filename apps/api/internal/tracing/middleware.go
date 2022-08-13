package tracing

import (
	"context"

	"contrib.rocks/apps/api/internal/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type contextKey int

const (
	traceNameContextKey contextKey = iota
)

func Middleware(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		originalCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(originalCtx)
		}()
		ctx := otel.GetTextMapPropagator().Extract(originalCtx, propagation.HeaderCarrier(c.Request.Header))
		ctx, span := Tracer().Start(ctx, "api.http")
		defer span.End()
		span.SetAttributes(attribute.String("app.environment", string(cfg.Env)))

		ctx = context.WithValue(ctx, traceNameContextKey, buildTraceName(cfg.ProjectID(), span.SpanContext().TraceID().String()))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TraceNameFromContext(c context.Context) string {
	if traceName, ok := c.Value(traceNameContextKey).(string); ok {
		return traceName
	}
	return ""
}

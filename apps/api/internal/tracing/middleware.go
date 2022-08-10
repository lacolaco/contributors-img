package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID := trace.SpanContextFromContext(ctx).TraceID()
		if !traceID.IsValid() {
			ctx, span := Tracer().Start(ctx, "api.http")
			defer span.End()
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		originalCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(originalCtx)
		}()
		ctx := otel.GetTextMapPropagator().Extract(originalCtx, propagation.HeaderCarrier(c.Request.Header))
		sc := trace.SpanContextFromContext(ctx)
		if !sc.HasSpanID() {
			var span trace.Span
			ctx, span = Tracer().Start(ctx, "api.http")
			defer span.End()
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

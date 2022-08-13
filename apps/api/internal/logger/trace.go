package logger

import (
	"context"

	"contrib.rocks/apps/api/internal/tracing"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
)

func withTracing(c context.Context) zap.Option {
	traceName := tracing.TraceNameFromContext(c)
	spanContext := trace.SpanContextFromContext(c)
	if traceName == "" || !spanContext.IsValid() {
		return zap.Fields()
	}
	return zap.Fields(
		zap.String(traceKey, traceName),
		zap.String(spanKey, spanContext.SpanID().String()),
		zap.Bool(traceSampledKey, spanContext.IsSampled()),
	)
}

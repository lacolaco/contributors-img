// Package cacheutil provides common utilities for cache operations
package cacheutil

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/internal/logger"
	"go.opentelemetry.io/otel/trace"
)

// LogCacheMiss logs cache miss events in a standardized format
func LogCacheMiss(ctx context.Context, cacheType string, key string) {
	logGroup := fmt.Sprintf("%s-cache-miss", cacheType)
	logger.LoggerFromContext(ctx).With(logger.LogGroup(logGroup)).Info(
		fmt.Sprintf("%s: %s", logGroup, key),
	)
}

// LogCacheSaveFailure records a cache write that did not land.
//
// Callers log rather than propagate: the response is already computed by this
// point, and failing the request would throw away work that succeeded — a
// rendered image, or contributor data that cost GitHub API quota to fetch. The
// cost of a lost write is latency on later requests, so it belongs in the logs,
// not in the status code.
func LogCacheSaveFailure(ctx context.Context, cacheType string, key string, err error) {
	trace.SpanFromContext(ctx).RecordError(err)
	logGroup := fmt.Sprintf("%s-cache-save-failure", cacheType)
	logger.LoggerFromContext(ctx).With(logger.LogGroup(logGroup)).Error(
		fmt.Sprintf("%s: %s: %s", logGroup, key, err),
	)
}

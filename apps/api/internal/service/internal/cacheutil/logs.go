// Package cacheutil provides common utilities for cache operations
package cacheutil

import (
	"context"
	"fmt"

	"contrib.rocks/apps/api/internal/logger"
)

// LogCacheMiss logs cache miss events in a standardized format
func LogCacheMiss(ctx context.Context, cacheType string, key string) {
	logGroup := fmt.Sprintf("%s-cache-miss", cacheType)
	logger.LoggerFromContext(ctx).With(logger.LogGroup(logGroup)).Info(
		fmt.Sprintf("%s: %s", logGroup, key),
	)
}

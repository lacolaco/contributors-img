package main

import (
	"context"

	"contrib.rocks/libs/goutils/config"
)

func SetEnvironment(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, config.EnvContextKey, env)
}

func GetEnvironment(ctx context.Context) string {
	return ctx.Value(config.EnvContextKey).(string)
}

package main

import (
	"context"

	"contrib.rocks/libs/goutils"
)

func SetEnvironment(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, goutils.EnvContextKey, env)
}

func GetEnvironment(ctx context.Context) string {
	return ctx.Value(goutils.EnvContextKey).(string)
}

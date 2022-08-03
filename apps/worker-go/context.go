package main

import "context"

type contextKey int

const (
	environmentKey contextKey = iota
)

func SetEnvironment(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, environmentKey, env)
}

func GetEnvironment(ctx context.Context) string {
	return ctx.Value(environmentKey).(string)
}

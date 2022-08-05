package goutils

import "os"

type Env string

const (
	EnvStaging    = Env("staging")
	EnvProduction = Env("production")
)

const (
	EnvContextKey = string("environment")
)

func GetEnv() Env {
	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		return EnvProduction
	case "staging":
		return EnvStaging
	default:
		return EnvStaging
	}
}

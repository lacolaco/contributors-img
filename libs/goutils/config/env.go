package config

import "os"

type Environment string

const (
	EnvStaging    = Environment("staging")
	EnvProduction = Environment("production")
)

const (
	EnvContextKey = string("environment")
)

func GetEnv() Environment {
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

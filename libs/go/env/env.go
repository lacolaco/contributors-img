package env

type Environment string

const (
	EnvDevelopment = Environment("development")
	EnvStaging     = Environment("staging")
	EnvProduction  = Environment("production")
)

// FromString app environment from string
func FromString(str string) Environment {
	switch str {
	case "production":
		return EnvProduction
	case "staging":
		return EnvStaging
	default:
		return EnvDevelopment
	}
}

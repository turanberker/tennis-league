package router

import "os"

type ServerConfig struct {
	Port           string
	AllowedOrigins string
	AppEnv         APP_ENV // "production" veya "development"
}

type APP_ENV string

const (
	APP_ENV_PRODUCTION  APP_ENV = "production"
	APP_END_DEVELOPMENT APP_ENV = "development"
)

func LoadServerConfig() *ServerConfig {

	env := getEnv("APP_ENV", "development")

	return &ServerConfig{
		Port:           getEnv("SERVER_PORT", "8500"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		AppEnv:         APP_ENV(env),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

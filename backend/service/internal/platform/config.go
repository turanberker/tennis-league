package platform

import "os"

type PostgreConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type RabbitConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	VHost    string
}

func LoadRabbitConfig() *RabbitConfig {
	return &RabbitConfig{
		User:     getEnv("RABBIT_USER", "admin"),
		Password: getEnv("RABBIT_PASSWORD", "admin123"),
		Host:     getEnv("RABBIT_HOST", "localhost"),
		Port:     getEnv("RABBIT_PORT", "5672"),
		VHost:    getEnv("RABBIT_VHOST", "/"),
	}
}

func LoadPostgresConfig() *PostgreConfig {
	return &PostgreConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "tennisleague"),
		Password: getEnv("POSTGRES_PASSWORD", "tennisleague"),
		Dbname:   getEnv("POSTGRES_DB", "tennisleague"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

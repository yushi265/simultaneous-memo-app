package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", buildDatabaseURL()),
		Environment: getEnv("GO_ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "dev")
	password := getEnv("DB_PASSWORD", "dev123")
	dbname := getEnv("DB_NAME", "notion_app")
	sslmode := getEnv("DB_SSLMODE", "disable")

	return "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
}
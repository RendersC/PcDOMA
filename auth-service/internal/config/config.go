package config

import (
	"os"
	"time"
)

type Config struct {
	Port           string
	DBDSN          string
	RedisURL       string
	JWTSecret      string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration
	Environment    string
}

func Load() *Config {
	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))

	return &Config{
		Port:          getEnv("PORT", "8081"),
		DBDSN:         getEnv("DB_DSN", ""),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret"),
		JWTAccessTTL:  accessTTL,
		JWTRefreshTTL: refreshTTL,
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

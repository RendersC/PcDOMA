package config

import "os"

type Config struct {
	Port              string
	DBDSN             string
	RedisURL          string
	CatalogServiceURL string
	PaymentServiceURL string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8083"),
		DBDSN:             getEnv("DB_DSN", ""),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/2"),
		CatalogServiceURL: getEnv("CATALOG_SERVICE_URL", "http://localhost:8082"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

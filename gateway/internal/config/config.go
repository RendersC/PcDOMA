package config

import "os"

type Config struct {
	Port                    string
	AuthServiceURL          string
	CatalogServiceURL       string
	BookingServiceURL       string
	PaymentServiceURL       string
	NotificationServiceURL  string
	RedisURL                string
	JWTSecret               string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		CatalogServiceURL:      getEnv("CATALOG_SERVICE_URL", "http://localhost:8082"),
		BookingServiceURL:      getEnv("BOOKING_SERVICE_URL", "http://localhost:8083"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8085"),
		RedisURL:               getEnv("REDIS_URL", "redis://localhost:6379/4"),
		JWTSecret:              getEnv("JWT_SECRET", "dev-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

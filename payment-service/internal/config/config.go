package config

import "os"

type Config struct {
	Port  string
	DBDSN string
}

func Load() *Config {
	return &Config{
		Port:  getEnv("PORT", "8084"),
		DBDSN: getEnv("DB_DSN", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

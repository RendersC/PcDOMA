package config

import "os"

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
	RedisURL string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8085"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "pcrentalnotifications"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/3"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

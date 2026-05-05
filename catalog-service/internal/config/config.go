package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
	RedisURL string
	CacheTTL time.Duration
}

func Load() *Config {
	ttlSec, _ := strconv.Atoi(getEnv("CACHE_TTL", "300"))
	return &Config{
		Port:     getEnv("PORT", "8082"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "pcrentalcatalog"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/1"),
		CacheTTL: time.Duration(ttlSec) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

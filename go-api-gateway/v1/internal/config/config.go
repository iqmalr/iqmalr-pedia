package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	AuthServiceURL  string
	UserServiceURL  string
	RedisURL        string
	RateLimit       int
	RateLimitWindow int
	JWTSecret       string
	Environment     string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))

	AppConfig = &Config{
		Port:            getEnv("PORT", "8080"),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		RateLimit:       rateLimit,
		RateLimitWindow: rateLimitWindow,
		JWTSecret:       getEnv("JWT_SECRET", "rahasia"),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

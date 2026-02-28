package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfigType struct {
	Port              string
	ProductServiceURL string
	AuthServiceURL    string
}

var AppConfig = AppConfigType{
	Port:              getEnv("PORT", "8085"),
	ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8084/api/v1"),
	AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://localhost:8082/api/v2"),
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

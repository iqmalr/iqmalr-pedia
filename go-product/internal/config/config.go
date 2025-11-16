package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfigType struct {
	Port             string
	VendorServiceURL string
}

var AppConfig = AppConfigType{
	Port:             getEnv("PORT", "8084"),
	VendorServiceURL: getEnv("VENDOR_SERVICE_URL", "http://localhost:8081/api/v1"),
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

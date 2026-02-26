package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DatabaseURL         string
	AuthServiceURL      string
	InternalAPIKey      string
	Environment         string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

var AppConfig *Config

func init() {
	LoadConfig()
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		Port:                getEnv("PORT", "8083"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:iqmalr@localhost:5432/vendors_db?sslmode=disable&timezone=Asia/Jakarta"),
		AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		InternalAPIKey:      getEnv("INTERNAL_API_KEY", "sangat-rahasia-juga-untuk-internal"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
	}

	log.Println("✅ Konfigurasi berhasil dimuat.")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// package config
//
// import (
//
//	"os"
//
//	"github.com/joho/godotenv"
//
// )
//
//	type Config struct {
//		Port           string
//		DatabaseURL    string
//		AuthServiceURL string
//		InternalAPIKey string
//		Environment    string
//	}
//
// var AppConfig *Config
//
//	func init() {
//		LoadConfig()
//	}
//
//	func LoadConfig() {
//		if err := godotenv.Load(); err != nil {
//			panic("Error loading .env file")
//		}
//
//		AppConfig = &Config{
//			Port:           getEnv("PORT", "8083"),
//			DatabaseURL:    getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/vendor_db?sslmode=disable"),
//			AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
//			InternalAPIKey: getEnv("INTERNAL_API_KEY", "sangat-rahasia-juga-untuk-internal"),
//			Environment:    getEnv("ENVIRONMENT", "development"),
//		}
//	}
//
//	func getEnv(key, defaultValue string) string {
//		if value := os.Getenv(key); value != "" {
//			return value
//		}
//		return defaultValue
//	}
package config

import (
	"log"
)

type Config struct {
	Port           string
	DatabaseURL    string
	AuthServiceURL string
	InternalAPIKey string
	Environment    string
}

var AppConfig *Config

// init() akan dijalankan otomatis saat package ini di-import
func init() {
	// Kita langsung isi konfigurasi dengan nilai hardcode
	AppConfig = &Config{
		Port: "8083",
		// Gunakan DSN dari file .env kamu yang lama
		DatabaseURL:    "postgres://postgres:iqmalr@localhost:5432/vendors_db?sslmode=disable&timezone=Asia/Jakarta",
		AuthServiceURL: "http://localhost:8080",
		InternalAPIKey: "sangat-rahasia-juga-untuk-internal",
		Environment:    "development",
	}

	log.Println("✅ Konfigurasi vendor (hardcode) berhasil dimuat.")
}

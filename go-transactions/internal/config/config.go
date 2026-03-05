package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfigType struct {
	Port                 string
	ProductServiceURL    string
	AuthServiceURL       string
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransEnv          string // "sandbox" or "production"
	PaymentProofUploadDir string
}

var AppConfig = AppConfigType{
	Port:                  getEnv("PORT", "8085"),
	ProductServiceURL:     getEnv("PRODUCT_SERVICE_URL", "http://localhost:8084/api/v1"),
	AuthServiceURL:        getEnv("AUTH_SERVICE_URL", "http://localhost:8082/api/v2"),
	MidtransServerKey:     getEnv("MIDTRANS_SERVER_KEY", ""),
	MidtransClientKey:     getEnv("MIDTRANS_CLIENT_KEY", ""),
	MidtransEnv:           getEnv("MIDTRANS_ENV", "sandbox"),
	PaymentProofUploadDir: getEnv("PAYMENT_PROOF_UPLOAD_DIR", "./uploads/payments"),
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

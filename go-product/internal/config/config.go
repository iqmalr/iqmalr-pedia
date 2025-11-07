package config

import (
	"os"
)

type AppConfigType struct {
	Port string
}

var AppConfig = AppConfigType{
	Port: getEnv("PORT", "8083"),
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

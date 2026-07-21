package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}
}

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

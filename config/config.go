package config

import (
	"os"
)

type Config struct {
	Port        string
	JWTSecret   string
	BotToken    string
	ProjectID   string
	FrontendURL string
}

var C Config

func Load() {
	C = Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		BotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		ProjectID:   os.Getenv("FIREBASE_PROJECT_ID"),
		FrontendURL: getEnv("FRONTEND_URL", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

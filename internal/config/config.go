package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is not set")
	}

	return Config{
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
	}, nil
}

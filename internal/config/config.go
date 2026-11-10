package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() (Config, error) {
	// Membaca file .env jika tersedia.
	// Di environment production, konfigurasi bisa berasal langsung dari OS.
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	return Config{
		DatabaseURL: databaseURL,
	}, nil
}

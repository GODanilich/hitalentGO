package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	GormLog     string
}

func Load() (Config, error) {

	_ = godotenv.Load()

	cfg := Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DB_URL"),
		GormLog:     getenv("GORM_LOG_LEVEL", "warn"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required (from .env or environment)")
	}
	return cfg, nil
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

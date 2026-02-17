package config

import (
	"os"
	"time"
)

type Config struct {
	Env         string
	HTTPAddr    string
	DatabaseURL string
	ReadTimeout time.Duration
}

func Load() Config {
	return Config{
		Env:         getenv("APP_ENV", "dev"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8081"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://app:app@localhost:5432/eventmesh?sslmode=disable"),
		ReadTimeout: 10 * time.Second,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

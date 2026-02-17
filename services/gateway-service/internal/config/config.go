package config

import (
	"os"
	"time"
)

type Config struct {
	Env         string
	HTTPAddr    string
	ReadTimeout time.Duration
}

func Load() Config {
	return Config{
		Env:         getenv("APP_ENV", "dev"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		ReadTimeout: 10 * time.Second,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

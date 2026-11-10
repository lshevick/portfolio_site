package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBPath             string
	ShopifyStoreURL    string
	ShopifyAccessToken string
}

func Load() Config {
	_ = godotenv.Load() // no-op if .env doesn't exist, e.g. in prod

	return Config{
		Port:   getenv("PORT", "localhost:8080"),
		DBPath: getenv("DB_PATH", "./data/app.db"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

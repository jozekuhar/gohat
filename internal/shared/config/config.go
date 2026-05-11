package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Host               string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	if cfg.Port, err = loadEnv("PORT"); err != nil {
		return nil, err
	}

	if cfg.Host, err = loadEnv("HOST"); err != nil {
		return nil, err
	}

	if cfg.DatabaseURL, err = loadEnv("DATABASE_URL"); err != nil {
		return nil, err
	}

	if cfg.GoogleClientID, err = loadEnv("GOOGLE_CLIENT_ID"); err != nil {
		return nil, err
	}

	if cfg.GoogleClientSecret, err = loadEnv("GOOGLE_CLIENT_SECRET"); err != nil {
		return nil, err
	}

	if cfg.GoogleRedirectURL, err = loadEnv("GOOGLE_REDIRECT_URL"); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("required env variable: %s", key)
	}
	return value, nil
}

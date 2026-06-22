package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Debug              bool
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

	if cfg.Debug, err = loadBoolEnv("DEBUG"); err != nil {
		return nil, err
	}

	if cfg.Port, err = loadStringEnv("PORT"); err != nil {
		return nil, err
	}

	if cfg.Host, err = loadStringEnv("HOST"); err != nil {
		return nil, err
	}

	if cfg.DatabaseURL, err = loadStringEnv("DATABASE_URL"); err != nil {
		return nil, err
	}

	if cfg.GoogleClientID, err = loadStringEnv("GOOGLE_CLIENT_ID"); err != nil {
		return nil, err
	}

	if cfg.GoogleClientSecret, err = loadStringEnv("GOOGLE_CLIENT_SECRET"); err != nil {
		return nil, err
	}

	if cfg.GoogleRedirectURL, err = loadStringEnv("GOOGLE_REDIRECT_URL"); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadStringEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("required env variable: %s", key)
	}
	return value, nil
}

func loadBoolEnv(key string) (bool, error) {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("required env variable: %s", key)
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}

	return value, nil
}

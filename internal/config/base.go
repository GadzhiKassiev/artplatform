package config

import "os"

type BaseConfig struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func loadBase(defaultPort string) BaseConfig {
	return BaseConfig{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
		Port:        getEnv("PORT", defaultPort),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

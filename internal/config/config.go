package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AuthURL     string
	CourseURL   string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
		AuthURL:     getEnv("AUTH_URL", "http://localhost:8001"),
		CourseURL:   getEnv("COURSE_URL", "http://localhost:8002"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

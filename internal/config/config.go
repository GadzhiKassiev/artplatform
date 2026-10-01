package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	ActivityURL string
	AuthURL     string
	CourseURL   string
	MediaURL    string
	PaymentURL  string
	PurchaseURL string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
		AuthURL:     getEnv("AUTH_URL", "http://localhost:8001"),
		CourseURL:   getEnv("COURSE_URL", "http://localhost:8002"),
		MediaURL:    getEnv("MEDIA_URL", "http://localhost:8003"),
		PaymentURL:  getEnv("PAYMENT_URL", "http://localhost:8004"),
		PurchaseURL: getEnv("PURCHASE_URL", "http://localhost:8005"),
		ActivityURL: getEnv("ACTIVITY_URL", "http://localhost:8006"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

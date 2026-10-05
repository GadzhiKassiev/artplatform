package config

type GatewayConfig struct {
	Port        string
	ActivityURL string
	AuthURL     string
	CourseURL   string
	MediaURL    string
	PaymentURL  string
	PurchaseURL string
}

func LoadGatewayConfig() GatewayConfig {
	return GatewayConfig{
		Port:        getEnv("PORT", "8000"),
		AuthURL:     getEnv("AUTH_URL", "http://localhost:8001"),
		CourseURL:   getEnv("COURSE_URL", "http://localhost:8002"),
		MediaURL:    getEnv("MEDIA_URL", "http://localhost:8003"),
		PaymentURL:  getEnv("PAYMENT_URL", "http://localhost:8004"),
		PurchaseURL: getEnv("PURCHASE_URL", "http://localhost:8005"),
		ActivityURL: getEnv("ACTIVITY_URL", "http://localhost:8006"),
	}
}

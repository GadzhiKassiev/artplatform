package config

import "strings"

type PurchaseConfig struct {
	BaseConfig
	Port         string
	KafkaBrokers []string
	PaymentURL   string
}

func LoadPurchaseConfig() PurchaseConfig {
	return PurchaseConfig{
		BaseConfig:   loadBase(),
		Port:         getEnv("PORT", "8005"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
		PaymentURL:   getEnv("PAYMENT_URL", "payment:8004"),
	}
}

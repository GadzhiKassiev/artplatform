package config

import "strings"

type PurchaseConfig struct {
	BaseConfig
	KafkaBrokers []string
	PaymentURL   string
}

func LoadPurchaseConfig() PurchaseConfig {
	return PurchaseConfig{
		BaseConfig:   loadBase("8005"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
		PaymentURL:   getEnv("PAYMENT_URL", "payment:8004"),
	}
}

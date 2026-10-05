package config

import "strings"

type PaymentConfig struct {
	BaseConfig
	Port         string
	KafkaBrokers []string
}

func LoadPaymentConfig() PaymentConfig {
	return PaymentConfig{
		BaseConfig:   loadBase(),
		Port:         getEnv("PORT", "8004"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
	}
}

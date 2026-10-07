package config

import "strings"

type PaymentConfig struct {
	BaseConfig
	KafkaBrokers []string
}

func LoadPaymentConfig() PaymentConfig {
	return PaymentConfig{
		BaseConfig:   loadBase("8004"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
	}
}

package config

import (
	"strings"
)

type ActivityConfig struct {
	Port         string
	RedisAddr    string
	KafkaBrokers []string
}

func LoadActivityConfig() ActivityConfig {
	return ActivityConfig{
		Port:         getEnv("PORT", "8006"),
		RedisAddr:    getEnv("REDIS_ADDR", "redis:6379"),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
	}
}

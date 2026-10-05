package config

import (
	"strings"
)

type MediaConfig struct {
	BaseConfig
	Port         string
	MinIO        MinIOConfig
	KafkaBrokers []string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func LoadMediaConfig() MediaConfig {
	return MediaConfig{
		BaseConfig: loadBase(),
		Port:       getEnv("PORT", "8003"),
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "artplatform-media"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
	}
}

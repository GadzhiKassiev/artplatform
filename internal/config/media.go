package config

import (
	"strings"
	"time"
)

type MediaConfig struct {
	BaseConfig
	MinIO             MinIOConfig
	KafkaBrokers      []string
	ProcessingTimeout time.Duration
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func LoadMediaConfig() MediaConfig {
	timeout, _ := time.ParseDuration(getEnv("PROCESSING_TIMEOUT", "30s"))
	return MediaConfig{
		BaseConfig: loadBase("8003"),
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "artplatform-media"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		ProcessingTimeout: timeout,
		KafkaBrokers:      strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
	}
}

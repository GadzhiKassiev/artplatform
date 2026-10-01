package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/media"
	"artplatform/backend/internal/service/media/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	"artplatform/backend/internal/transport/kafka"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	minioClient, err := storage.NewMinIO(ctx, storage.MinIOConfig{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET"),
		UseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
	})
	if err != nil {
		logging.Fatal("failed to connect to minio", logging.NewKV("error", err))
	}

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	logging.Info(ctx, "waiting for kafka topic", logging.NewKV("topic", media.TopicMediaUploaded))
	if err := kafka.WaitForKafka(brokers, media.TopicMediaUploaded, 30*time.Second); err != nil {
		logging.Fatal("failed to ensure kafka topic", logging.NewKV("error", err))
	}
	logging.Info(ctx, "kafka topic ready", logging.NewKV("topic", media.TopicMediaUploaded))
	st := storage.NewPostgres(pool)
	worker := media.NewWorker(st, minioClient, brokers)
	defer worker.Close()

	logging.Info(ctx, "media-worker starting")
	if err := worker.Run(ctx); err != nil {
		logging.Fatal("worker failed", logging.NewKV("error", err))
	}
	logging.Info(ctx, "media-worker stopped")
}

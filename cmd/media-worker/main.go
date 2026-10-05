package main

import (
	"context"
	"os/signal"
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

	cfg := config.LoadMediaConfig()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	minioClient, err := storage.NewMinIO(ctx, storage.MinIOConfig{
		Endpoint:  cfg.MinIO.Endpoint,
		AccessKey: cfg.MinIO.AccessKey,
		SecretKey: cfg.MinIO.SecretKey,
		Bucket:    cfg.MinIO.Bucket,
		UseSSL:    cfg.MinIO.UseSSL,
	})
	if err != nil {
		logging.Fatal("failed to connect to minio", logging.NewKV("error", err))
	}

	logging.Info(ctx, "waiting for kafka topic", logging.NewKV("topic", media.TopicMediaUploaded))
	if err := kafka.WaitForKafka(cfg.KafkaBrokers, media.TopicMediaUploaded, 30*time.Second); err != nil {
		logging.Fatal("failed to ensure kafka topic", logging.NewKV("error", err))
	}
	logging.Info(ctx, "kafka topic ready", logging.NewKV("topic", media.TopicMediaUploaded))
	st := storage.NewPostgres(pool)
	worker := media.NewWorker(st, minioClient, cfg.KafkaBrokers)
	defer worker.Close()

	logging.Info(ctx, "media-worker starting")
	if err := worker.Run(ctx); err != nil {
		logging.Fatal("worker failed", logging.NewKV("error", err))
	}
	logging.Info(ctx, "media-worker stopped")
}

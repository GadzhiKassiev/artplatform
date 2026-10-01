package main

import (
	"context"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/media"
	"artplatform/backend/internal/service/media/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	grpctransport "artplatform/backend/internal/transport/grpc"
	"artplatform/backend/internal/transport/kafka"
	mediapb "artplatform/backend/proto/media"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)
	ctx := context.Background()

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
	producer := kafka.NewProducer(brokers)
	defer producer.Close()

	st := storage.NewPostgres(pool)
	svc := media.New(st, minioClient, producer)
	grpcServer := media.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	mediapb.RegisterMediaServiceServer(s, grpcServer)

	logging.Info(ctx, "media service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

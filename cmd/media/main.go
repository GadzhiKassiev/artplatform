package main

import (
	"context"
	"net"

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

	mediaCfg := config.LoadMediaConfig()

	pool, err := pgstorage.New(ctx, mediaCfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	minioClient, err := storage.NewMinIO(ctx, storage.MinIOConfig{
		Endpoint:  mediaCfg.MinIO.Endpoint,
		AccessKey: mediaCfg.MinIO.AccessKey,
		SecretKey: mediaCfg.MinIO.SecretKey,
		Bucket:    mediaCfg.MinIO.Bucket,
		UseSSL:    mediaCfg.MinIO.UseSSL,
	})
	if err != nil {
		logging.Fatal("failed to connect to minio", logging.NewKV("error", err))
	}

	producer := kafka.NewProducer(mediaCfg.KafkaBrokers)
	defer producer.Close()

	st := storage.NewPostgres(pool)
	svc := media.New(st, minioClient, producer)
	grpcServer := media.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+mediaCfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	mediapb.RegisterMediaServiceServer(s, grpcServer)

	logging.Info(ctx, "media service listening", logging.NewKV("port", mediaCfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

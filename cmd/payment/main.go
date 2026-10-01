package main

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/payment"
	"artplatform/backend/internal/service/payment/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	grpctransport "artplatform/backend/internal/transport/grpc"
	"artplatform/backend/internal/transport/kafka"
	paymentpb "artplatform/backend/proto/payment"
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

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	if err := kafka.WaitForKafka(brokers, payment.TopicPaymentSucceeded, 30*time.Second); err != nil {
		logging.Fatal("failed to ensure kafka topic", logging.NewKV("error", err))
	}

	producer := kafka.NewProducer(brokers)
	defer producer.Close()

	st := storage.NewPostgres(pool)
	svc := payment.New(st, producer)
	grpcServer := payment.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	paymentpb.RegisterPaymentServiceServer(s, grpcServer)

	logging.Info(ctx, "payment service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

package main

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/purchase"
	"artplatform/backend/internal/service/purchase/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	grpctransport "artplatform/backend/internal/transport/grpc"
	"artplatform/backend/internal/transport/kafka"
	purchasepb "artplatform/backend/proto/purchase"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)
	ctx := context.Background()

	cfg := config.LoadPurchaseConfig()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	if err := kafka.WaitForKafka(cfg.KafkaBrokers, purchase.TopicPurchaseCreated, 30*time.Second); err != nil {
		logging.Fatal("failed to ensure kafka topic", logging.NewKV("error", err))
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	paymentClient, paymentConn, err := grpctransport.NewPaymentClient(cfg.PaymentURL)
	if err != nil {
		logging.Fatal("failed to connect to payment", logging.NewKV("error", err))
	}
	defer paymentConn.Close()

	st := storage.NewPostgres(pool)
	svc := purchase.New(st, producer, paymentClient)
	grpcServer := purchase.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	purchasepb.RegisterPurchaseServiceServer(s, grpcServer)

	logging.Info(ctx, "purchase service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

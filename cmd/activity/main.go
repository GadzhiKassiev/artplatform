package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/activity"
	"artplatform/backend/internal/service/activity/storage"
	redisstorage "artplatform/backend/internal/storage/redis"
	grpctransport "artplatform/backend/internal/transport/grpc"
	"artplatform/backend/internal/transport/kafka"
	activitypb "artplatform/backend/proto/activity"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load()

	redisClient, err := redisstorage.New(ctx, os.Getenv("REDIS_ADDR"))
	if err != nil {
		logging.Fatal("failed to connect to redis", logging.NewKV("error", err))
	}
	defer redisClient.Close()

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	for _, topic := range []string{activity.TopicCourseViewed, activity.TopicPurchaseCreated} {
		if err := kafka.WaitForKafka(brokers, topic, 30*time.Second); err != nil {
			logging.Fatal("failed to ensure topic", logging.NewKV("topic", topic), logging.NewKV("error", err))
		}
	}

	st := storage.NewRedis(redisClient)
	svc := activity.New(st)
	grpcServer := activity.NewGRPCServer(svc)

	go func() {
		if err := svc.StartConsumers(ctx, brokers); err != nil {
			logging.ErrorE(ctx, "consumers failed", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	activitypb.RegisterActivityServiceServer(s, grpcServer)

	logging.Info(ctx, "activity service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

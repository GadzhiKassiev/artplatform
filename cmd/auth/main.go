package main

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/jwt"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/auth"
	"artplatform/backend/internal/service/auth/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	grpctransport "artplatform/backend/internal/transport/grpc"
	authpb "artplatform/backend/proto/auth"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)
	ctx := context.Background()

	cfg := config.LoadAuthConfig()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	st := storage.NewPostgres(pool)
	tokenizer := jwt.New(cfg.JWTSecret, 24*time.Hour)
	svc := auth.New(st, tokenizer)
	grpcServer := auth.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.ValidationInterceptor()),
	)
	authpb.RegisterAuthServiceServer(s, grpcServer)

	logging.Info(ctx, "auth service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

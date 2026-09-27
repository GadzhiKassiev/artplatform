package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/pkg/jwt"
	"artplatform/backend/internal/service/auth"
	"artplatform/backend/internal/service/auth/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	authpb "artplatform/backend/proto/auth"
)

func main() {
	logging.Init()

	cfg := config.Load()

	ctx := context.Background()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	st := storage.NewPostgres(pool)
	tokenizer := jwt.New(cfg.JWTSecret, 24*time.Hour)
	svc := auth.New(st, tokenizer)
	grpcServer := auth.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, grpcServer)

	log.Printf("auth service listening on :%s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/course"
	"artplatform/backend/internal/service/course/storage"
	pgstorage "artplatform/backend/internal/storage/postgres"
	coursepb "artplatform/backend/proto/course"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)

	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Fatal("failed to connect to postgres", logging.NewKV("error", err))
	}
	defer pool.Close()

	st := storage.NewPostgres(pool)
	svc := course.New(st)
	grpcServer := course.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logging.Fatal("failed to listen", logging.NewKV("error", err))
	}

	s := grpc.NewServer()
	coursepb.RegisterCourseServiceServer(s, grpcServer)

	logging.Info(ctx, "course service listening", logging.NewKV("port", cfg.Port))
	if err := s.Serve(lis); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

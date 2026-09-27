package main

import (
	"context"
	"log"
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
	logging.Init()

	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgstorage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	st := storage.NewPostgres(pool)
	svc := course.New(st)
	grpcServer := course.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	coursepb.RegisterCourseServiceServer(s, grpcServer)

	log.Printf("course service listening on :%s", cfg.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

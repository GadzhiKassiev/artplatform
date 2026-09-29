package main

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"artplatform/backend/api/graphql"
	"artplatform/backend/api/graphql/generated"
	"artplatform/backend/internal/config"
	"artplatform/backend/internal/logging"
	grpctransport "artplatform/backend/internal/transport/grpc"
)

func main() {
	logging.Init(logging.LevelInfo, logging.FormatJSON)

	cfg := config.Load()
	ctx := context.Background()

	authClient, authConn, err := grpctransport.NewAuthClient(cfg.AuthURL)
	if err != nil {
		logging.Fatal("failed to connect to auth", logging.NewKV("error", err))
	}
	defer authConn.Close()

	courseClient, courseConn, err := grpctransport.NewCourseClient(cfg.CourseURL)
	if err != nil {
		logging.Fatal("failed to connect to course", logging.NewKV("error", err))
	}
	defer courseConn.Close()

	resolver := &graphql.Resolver{
		AuthClient:   authClient,
		CourseClient: courseClient,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))
	srv.SetErrorPresenter(graphql.ErrorPresenter)

	authMW := &graphql.AuthMiddleware{AuthClient: authClient}

	mux := http.NewServeMux()
	mux.Handle("/query", authMW.Middleware(srv))
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	logging.Info(ctx, "gateway listening", logging.NewKV("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

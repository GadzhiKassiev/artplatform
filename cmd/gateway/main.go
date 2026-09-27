package main

import (
	"log"
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
	logging.Init()

	cfg := config.Load()

	authClient, authConn, err := grpctransport.NewAuthClient(cfg.AuthURL)
	if err != nil {
		log.Fatalf("failed to connect to auth: %v", err)
	}
	defer authConn.Close()

	courseClient, courseConn, err := grpctransport.NewCourseClient(cfg.CourseURL)
	if err != nil {
		log.Fatalf("failed to connect to course: %v", err)
	}
	defer courseConn.Close()

	resolver := &graphql.Resolver{
		AuthClient:   authClient,
		CourseClient: courseClient,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	authMW := &graphql.AuthMiddleware{AuthClient: authClient}

	mux := http.NewServeMux()
	mux.Handle("/query", authMW.Middleware(srv))
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	log.Printf("gateway listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

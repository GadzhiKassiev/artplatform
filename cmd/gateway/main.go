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
	ctx := context.Background()

	cfg := config.LoadGatewayConfig()

	authClient, err := grpctransport.NewAuthClient(cfg.AuthURL)
	if err != nil {
		logging.Fatal("failed to connect to auth", logging.NewKV("error", err))
	}
	defer authClient.Close()

	courseClient, err := grpctransport.NewCourseClient(cfg.CourseURL)
	if err != nil {
		logging.Fatal("failed to connect to course", logging.NewKV("error", err))
	}
	defer courseClient.Close()

	mediaClient, err := grpctransport.NewMediaClient(cfg.MediaURL)
	if err != nil {
		logging.Fatal("failed to connect to media", logging.NewKV("error", err))
	}
	defer mediaClient.Close()

	paymentClient, err := grpctransport.NewPaymentClient(cfg.PaymentURL)
	if err != nil {
		logging.Fatal("failed to connect to payment", logging.NewKV("error", err))
	}
	defer paymentClient.Close()

	purchaseClient, err := grpctransport.NewPurchaseClient(cfg.PurchaseURL)
	if err != nil {
		logging.Fatal("failed to connect to purchase", logging.NewKV("error", err))
	}
	defer purchaseClient.Close()

	activityClient, err := grpctransport.NewActivityClient(cfg.ActivityURL)
	if err != nil {
		logging.Fatal("failed to connect to activity", logging.NewKV("error", err))
	}
	defer activityClient.Close()

	resolver := &graphql.Resolver{
		AuthClient:     authClient.Client,
		CourseClient:   courseClient.Client,
		MediaClient:    mediaClient.Client,
		PaymentClient:  paymentClient.Client,
		PurchaseClient: purchaseClient.Client,
		ActivityClient: activityClient.Client,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))
	srv.SetErrorPresenter(graphql.ErrorPresenter)

	authMW := &graphql.AuthMiddleware{AuthClient: authClient.Client}

	mux := http.NewServeMux()
	mux.Handle("/query", authMW.Middleware(srv))
	mux.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	logging.Info(ctx, "gateway listening", logging.NewKV("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		logging.Fatal("failed to serve", logging.NewKV("error", err))
	}
}

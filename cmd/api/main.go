package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/FarhanDH/sidingin-rating-service/config"
	httpHandler "github.com/FarhanDH/sidingin-rating-service/internal/handler"
	grpcHandler "github.com/FarhanDH/sidingin-rating-service/internal/handler/grpc"
	"github.com/FarhanDH/sidingin-rating-service/internal/platform/auth"
	"github.com/FarhanDH/sidingin-rating-service/internal/repository"
	"github.com/FarhanDH/sidingin-rating-service/internal/service"
	"github.com/FarhanDH/sidingin-rating-service/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	handler := httpHandler.NewHandler()
	port := ":" + cfg.Server.RestPort
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to open db %v:", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to pin db %v:", err)
	}

	jwt, err := auth.NewJwtVerifier(cfg.JWT.Secret, cfg.JWT.Issuer)
	if err != nil {
		log.Fatalf("failed to load jwt verifier: %v", err)
	}

	repo := repository.NewPostgresRepository(db)
	svc := service.NewRatingService(repo)
	grpcRatingHandler := grpcHandler.NewGRPCServer(svc, jwt)
	grpcPort := ":" + cfg.Server.GRPCPort
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRatingServiceServer(grpcServer, grpcRatingHandler)

	log.Printf("gRPC server running on port %s", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC Server: %v", err)
	}

	// HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	if err := http.ListenAndServe(port, mux); err != nil {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}

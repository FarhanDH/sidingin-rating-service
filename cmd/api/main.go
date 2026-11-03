package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FarhanDH/sidingin-rating-service/config"
	httpHandler "github.com/FarhanDH/sidingin-rating-service/internal/handler"
	"github.com/FarhanDH/sidingin-rating-service/internal/handler/event"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	jwt, err := auth.NewJwtVerifier(cfg.JWT.Secret, cfg.JWT.Issuer)
	if err != nil {
		logger.Error("failed to load jwt verifier", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo := repository.NewPostgresRepository(db)
	svc := service.NewRatingService(repo)

	// HTTP handler
	handler := httpHandler.NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	// RabbitMQ consumer
	rabbitConsumer := event.NewRabbitMQHandler(&cfg.RabbitMQ, svc, logger)
	if err := rabbitConsumer.Connect(context.Background()); err != nil {
		logger.Error("failed to connect to RabbitMQ", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer rabbitConsumer.Close()

	// Start consumer in background
	go func() {
		if err := rabbitConsumer.Consume(context.Background()); err != nil && err != context.Canceled {
			logger.Error("RabbitMQ consumer stopped unexpectedly", slog.String("error", err.Error()))
		}
	}()

	// gRPC server
	grpcRatingHandler := grpcHandler.NewGRPCServer(svc, jwt)
	grpcPort := ":" + cfg.Server.GRPCPort
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Error("failed to listen on gRPC port", slog.String("port", grpcPort), slog.String("error", err.Error()))
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRatingServiceServer(grpcServer, grpcRatingHandler)

	// Run gRPC in background goroutine
	go func() {
		logger.Info("gRPC server running", slog.String("port", grpcPort))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server stopped unexpectedly", slog.String("error", err.Error()))
		}
	}()

	// HTTP server
	httpPort := ":" + cfg.Server.RestPort
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}
	go func() {
		logger.Info("HTTP server running", slog.String("port", httpPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server stopped unexpectedly", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop gRPC
	grpcServer.GracefulStop()

	// Stop HTTP
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", slog.String("error", err.Error()))
	}

	// Close RabbitMQ and DB
	rabbitConsumer.Close()
	db.Close()

	logger.Info("server stopped")
}

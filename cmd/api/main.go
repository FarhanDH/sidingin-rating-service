package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/FarhanDH/sidingin-rating-service/config"
	httpHandler "github.com/FarhanDH/sidingin-rating-service/internal/handler"
)

func main() {
	cfg := config.Load()
	handler := httpHandler.NewHandler()
	port := ":" + cfg.Server.Port
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	if err := http.ListenAndServe(port, mux); err != nil {
		logger.Error("server failed to start: %v", "error", err)
		os.Exit(1)
	}
}

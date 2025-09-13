package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"google-auth-demo/backend/internal/config"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/oauth/google"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"

	"google-auth-demo/backend/internal/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service got error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Printf("DEBUG: HttpServer Port=%d\n", cfg.HttpServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := logger.Init(cfg.Logger); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	return start(ctx, cfg)
}

func start(ctx context.Context, cfg *config.Config) error {
	oauthGoogle := google.New(cfg.GoogleAuth)
	repository := repo.NewMockRepo(cfg.Repo)
	svc := service.New(cfg.Service, oauthGoogle, repository)
	httpServer := httpserver.New(cfg.HttpServer, svc)

	return httpServer.Run(ctx)
}

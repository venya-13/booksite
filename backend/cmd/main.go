package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"

	"google-auth-demo/backend/internal/logger"
	googloauth "google-auth-demo/backend/internal/oauth/google"
)

func main() {
	if err := run(); err != nil {
		slog.Error("service got error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfigFromEnvs()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := logger.Init(cfg.Logger); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	return start(ctx, cfg)
}

func start(ctx context.Context, cfg *Config) error {
	oauthGoogle := googloauth.New(cfg.GoogleAuth)

	repository := repo.NewMockRepo(cfg.Repo)

	svc := service.New(cfg.Service, oauthGoogle, repository)

	httpServer := httpserver.New(cfg.HttpServer, svc)

	return httpServer.Run(ctx)
}

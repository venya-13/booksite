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
	"google-auth-demo/backend/internal/logger"
	"google-auth-demo/backend/internal/oauth/google"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"
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

	if err := logger.Init(logger.Config{
		Level: cfg.Logger.Level,
	}); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	return start(ctx, cfg)
}

func start(ctx context.Context, cfg *config.Config) error {

	googleCfg := google.Config{
		ClientID:        cfg.GoogleAuth.ClientID,
		ClientSecret:    cfg.GoogleAuth.ClientSecret,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
	}

	oauthGoogle := google.New(googleCfg)
	repository := repo.NewMockRepo(repo.Config{})

	svc := service.New(service.Config{
		FrontendURL: cfg.HttpServer.FrontendURL,
	}, oauthGoogle, repository)

	httpServerCfg := httpserver.Config{
		Port:            cfg.HttpServer.Port,
		FrontendURL:     cfg.HttpServer.FrontendURL,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
	}

	httpServer := httpserver.New(httpServerCfg, svc)

	fmt.Println("DEBUG: HttpServer Port:", cfg.HttpServer.Port)
	fmt.Println("DEBUG: Frontend URL:", cfg.HttpServer.FrontendURL)
	fmt.Println("DEBUG: Redirect Base URL:", cfg.HttpServer.RedirectBaseURL)
	fmt.Println("DEBUG: Google ClientID:", cfg.GoogleAuth.ClientID)
	fmt.Println("DEBUG: Google ClientSecret:", cfg.GoogleAuth.ClientSecret != "")
	fmt.Println("DEBUG: Logger Level:", cfg.Logger.Level)

	return httpServer.Run(ctx)
}

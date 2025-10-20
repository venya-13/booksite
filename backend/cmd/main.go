package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google-auth-demo/backend/internal/config"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/jwt"
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

	if err := logger.Init(logger.Config{
		Level: cfg.Logger.Level,
		JSON:  cfg.Logger.JSON,
	}); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	slog.Info("Logger initialized", slog.String("level", cfg.Logger.Level))

	jwt.SetJWTKey(cfg.JWT.Secret)
	slog.Info("JWT secret set")

	slog.Info("Starting application",
		slog.Int("port", cfg.HttpServer.Port),
		slog.String("frontend_url", cfg.HttpServer.FrontendURL),
		slog.String("redirect_base", cfg.HttpServer.RedirectBaseURL),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return start(ctx, cfg)
}

func start(ctx context.Context, cfg *config.Config) error {

	googleCfg := google.Config{
		ClientID:        cfg.GoogleAuth.ClientID,
		ClientSecret:    cfg.GoogleAuth.ClientSecret,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
		Port:            cfg.HttpServer.Port,
	}

	oauthGoogle := google.New(googleCfg)

	repository, err := repo.NewPostgresRepo(repo.PostgresConfig{
		DSN: cfg.Database.DSN,
	})

	if err != nil {
		return fmt.Errorf("init repo: %w", err)
	}

	jwtTTL := time.Duration(cfg.JWT.TTL) * time.Minute
	refreshTTL := time.Duration(cfg.JWT.RefreshTTL) * time.Minute

	svc := service.New(
		cfg.HttpServer.FrontendURL,
		oauthGoogle,
		repository,
		jwtTTL,
		refreshTTL,
	)

	httpServerCfg := httpserver.Config{
		Port:            cfg.HttpServer.Port,
		FrontendURL:     cfg.HttpServer.FrontendURL,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
	}

	httpServer := httpserver.New(httpServerCfg, svc)

	return httpServer.Run(ctx)
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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

	fmt.Fprintln(os.Stderr, "[debug] config loaded, about to init logger")

	slog.Info("Config values",
		"port", cfg.HttpServer.Port,
		"frontend_url", cfg.HttpServer.FrontendURL,
		"redirect_base_url", cfg.HttpServer.RedirectBaseURL,
		"jwt_secret", cfg.JWT.Secret,
	)

	if err := logger.Init(logger.Config{
		Level: cfg.Logger.Level,
		JSON:  cfg.Logger.JSON,
	}); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	slog.Info("Logger initialized", slog.String("level", cfg.Logger.Level))

	slog.Info("Starting application",
		slog.Int("port", cfg.HttpServer.Port),
		slog.String("frontend_url", cfg.HttpServer.FrontendURL),
		slog.String("redirect_base", cfg.HttpServer.RedirectBaseURL),
	)

	// On Windows, os.Interrupt is the reliable signal to listen for.
	// Using unsupported signals can result in surprising behavior.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Fprintln(os.Stderr, "[debug] signal context created; ctx.Err() =", ctx.Err())

	return start(ctx, cfg)
}

func start(ctx context.Context, cfg *config.Config) error {

	jwt.SetSecrets(cfg.JWT.Secret, cfg.JWT.Secret)

	fmt.Fprintln(os.Stderr, "[debug] starting: initializing oauth + repo")

	googleCfg := google.Config{
		ClientID:        cfg.GoogleAuth.ClientID,
		ClientSecret:    cfg.GoogleAuth.ClientSecret,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
		Port:            cfg.HttpServer.Port,
	}

	oauthGoogle := google.New(googleCfg)

	fmt.Fprintln(os.Stderr, "[debug] creating postgres repo...")
	repository, err := repo.NewPostgresRepo(repo.PostgresConfig{
		DSN: cfg.Database.DSN,
	})

	if err != nil {
		return fmt.Errorf("init repo: %w", err)
	}
	fmt.Fprintln(os.Stderr, "[debug] repo ready; ctx.Err() =", ctx.Err())

	jwtTTL := time.Duration(cfg.JWT.TTL) * time.Minute
	refreshTTL := time.Duration(cfg.JWT.RefreshTTL) * time.Minute

	svc := service.New(
		cfg.HttpServer.FrontendURL,
		oauthGoogle,
		repository,
		jwtTTL,
		refreshTTL,
	)

	fmt.Fprintln(os.Stderr, "[debug] building http server...")
	httpServerCfg := httpserver.Config{
		Port:            cfg.HttpServer.Port,
		FrontendURL:     cfg.HttpServer.FrontendURL,
		RedirectBaseURL: cfg.HttpServer.RedirectBaseURL,
	}

	httpServer := httpserver.New(httpServerCfg, svc)

	fmt.Fprintln(os.Stderr, "[debug] running http server; should listen on port", cfg.HttpServer.Port)
	return httpServer.Run(ctx)
}

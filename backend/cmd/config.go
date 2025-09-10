package main

import (
	"fmt"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/logger"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	googloauth "google-auth-demo/backend/internal/oauth/google"
)

type Config struct {
	Logger     logger.Config     `envprefix:"LOG_"`
	GoogleAuth googloauth.Config `envprefix:"GOOGLE_AUTH_"`
	Repo       repo.Config       `envprefix:"REPO_"`
	HttpServer httpserver.Config `envprefix:"HTTP_SERVER_"`
	Service    service.Config    `envprefix:"SERVICE_"`
}

func loadConfigFromEnvs() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "debug"
	}

	return &cfg, nil
}

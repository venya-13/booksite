package config

import (
	"fmt"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/logger"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"
	"os"

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

func Load() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting cwd: %w", err)
	}

	envPath := cwd + string(os.PathSeparator) + ".env"

	if err := godotenv.Load(envPath); err != nil {
		fmt.Println("Warning: .env file not found at", envPath)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}

	return &cfg, nil
}

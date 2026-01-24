package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServer HttpServerConfig `mapstructure:"httpserver"`
	GoogleAuth GoogleAuthConfig `mapstructure:"googleauth"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret" json:"secret" yaml:"secret"`
	TTL        int    `mapstructure:"ttl" json:"ttl" yaml:"ttl"` // in minutes
	RefreshTTL int    `mapstructure:"refresh_ttl" json:"refresh_ttl" yaml:"refresh_ttl"`
}

type HttpServerConfig struct {
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	FrontendURL     string `mapstructure:"frontend_url" json:"frontend_url" yaml:"frontend_url"`
	RedirectBaseURL string `mapstructure:"redirect_base_url" json:"redirect_base_url" yaml:"redirect_base_url"`
}

type GoogleAuthConfig struct {
	ClientID     string `mapstructure:"client_id" json:"client_id" yaml:"client_id"`
	ClientSecret string `mapstructure:"client_secret" json:"client_secret" yaml:"client_secret"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	JSON  bool   `mapstructure:"json" json:"json" yaml:"json"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Enable environment variables
	v.SetEnvPrefix("")
	v.AutomaticEnv()

	// Set environment variable mappings
	v.BindEnv("httpserver.port", "HTTPSERVER_PORT")
	v.BindEnv("httpserver.frontend_url", "HTTPSERVER_FRONTEND_URL")
	v.BindEnv("httpserver.redirect_base_url", "HTTPSERVER_REDIRECT_BASE_URL")
	v.BindEnv("database.dsn", "DATABASE_DSN")
	v.BindEnv("jwt.secret", "JWT_SECRET")
	v.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	v.BindEnv("jwt.ttl", "JWT_TTL")
	v.BindEnv("jwt.refresh_ttl", "JWT_REFRESH_TTL")
	v.BindEnv("logger.level", "LOGGER_LEVEL")
	v.BindEnv("logger.json", "LOGGER_JSON")
	v.BindEnv("googleauth.client_id", "GOOGLE_CLIENT_ID")
	v.BindEnv("googleauth.client_secret", "GOOGLE_CLIENT_SECRET")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./..")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		slog.Warn("Config file not found, trying JSON", "error", err)
		v.SetConfigType("json")
		if err := v.ReadInConfig(); err != nil {
			slog.Warn("Config JSON file not found, using environment variables or defaults", "error", err)
		}
	} else {
		slog.Info("Config file loaded", "file", v.ConfigFileUsed())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

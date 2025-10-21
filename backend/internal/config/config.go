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

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./..")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		slog.Warn("Config file not found, trying JSON", "error", err)
		v.SetConfigType("json")
		if err := v.ReadInConfig(); err != nil {
			slog.Warn("Config JSON file not found, using defaults", "error", err)
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

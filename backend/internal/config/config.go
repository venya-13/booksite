package config

import (
	"fmt"

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
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	TTL    int    `mapstructure:"ttl" json:"ttl" yaml:"ttl"` // in minutes
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

	// try json first
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		// try yaml next
		v.SetConfigType("yaml")
		if err2 := v.ReadInConfig(); err2 != nil {
			fmt.Println("Warning: config.json and config.yaml not found, using defaults")
		}
	}

	// defaults
	v.SetDefault("HttpServer.Port", 8080)
	v.SetDefault("Logger.Level", "info")
	viper.SetDefault("JWT.Secret", "super-secret-key")
	viper.SetDefault("JWT.TTL", 60) // 60 minutes

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}

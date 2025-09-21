package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HttpServer struct {
		Port            int    `mapstructure:"port"`
		FrontendURL     string `mapstructure:"frontend_url"`
		RedirectBaseURL string `mapstructure:"redirect_base_url"`
	} `mapstructure:"httpserver"`

	GoogleAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
	} `mapstructure:"googleauth"`

	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`

	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`
}

func Load() (*Config, error) {
	viper.SetConfigFile("config.yaml") // location of config file
	viper.AutomaticEnv()               // yaml

	viper.SetDefault("HttpServer.Port", 8080)
	viper.SetDefault("Logger.Level", "info")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: config.yaml not found, using defaults and env")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}

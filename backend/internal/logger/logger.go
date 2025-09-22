package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Level string `env:"LEVEL" envDefault:"info" json:"level" yaml:"level"`
	JSON  bool   `env:"JSON" json:"json" yaml:"json"`
}

const DefaultLevel = "info"

func Init(config Config) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(config.Level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", config.Level, err)
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}

	if config.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

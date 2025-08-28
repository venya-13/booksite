package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Level string `env:"LEVEL" envDefault:"info"`
	JSON  bool   `env:"JSON"`
}

func Init(config Config) error {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(config.Level)); err != nil {
		return fmt.Errorf("invalid level: %w", err)
	}

	var logHandler slog.Handler
	if config.JSON {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	}

	slog.SetDefault(slog.New(logHandler))

	return nil
}

package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"google-auth-demo/backend/internal/service"
)

type (
	Server struct {
		svc *service.Service
		s   *http.Server
	}

	Config struct {
		Port            int    `env:"PORT"`
		RedirectBaseURL string `env:"REDIRECT_BASE_URL"`
		FrontendURL     string `env:"FRONTEND_URL"`
	}
)

func New(config Config, svc *service.Service) *Server {

	srv := Server{
		svc: svc,
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: srv.createMux(),
	}
	srv.s = &httpServer

	return &srv
}

func (s *Server) createMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/oauth2callback", s.handleCallback)

	return mux
}

func (s *Server) Run(ctx context.Context) error {
	slog.Info("server starting", slog.String("addr", s.s.Addr))

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.s.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutting down server", slog.String("error", err.Error()))
		}
	}()

	if err := s.s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server close: %w", err)
		}
	}

	return nil
}

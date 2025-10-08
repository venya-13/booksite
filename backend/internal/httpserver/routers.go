package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"google-auth-demo/backend/internal/httpserver/middleware"
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
	mux.HandleFunc("/refresh", s.handleRefresh)
	mux.HandleFunc("/api/auth/refresh", s.handleJWTRefresh)
	mux.Handle("/api/google-profile", middleware.AuthMiddleware(http.HandlerFunc(s.handleGoogleProfile)))
	mux.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(s.handleProtected)))

	return mux
}

func (s *Server) Run(ctx context.Context) error {
	slog.Info("HTTP server starting", slog.String("address", s.s.Addr))

	go func() {
		<-ctx.Done()
		slog.Warn("Shutdown signal received, stopping server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.s.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error during server shutdown", slog.String("error", err.Error()))
		} else {
			slog.Info("Server shutdown completed cleanly")
		}
	}()

	if err := s.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server error: %w", err)
	}

	return nil
}

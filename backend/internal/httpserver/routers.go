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

func (s *Server) createMux() http.Handler {
	mux := http.NewServeMux()

	// --- AUTH ---
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/oauth2callback", s.handleCallback)
	mux.HandleFunc("/refresh", s.handleRefresh)
	mux.HandleFunc("/api/auth/refresh", s.handleJWTRefresh)

	// --- GOOGLE PROFILE + PROTECTED ROUTES ---
	mux.Handle("/api/google-profile", middleware.AuthMiddleware(http.HandlerFunc(s.handleGoogleProfile)))
	mux.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(s.handleProtected)))

	// --- CATEGORIES ---
	mux.HandleFunc("/api/categories", s.handleGetCategories)
	mux.Handle("/api/categories/create", middleware.AuthMiddleware(http.HandlerFunc(s.handleCreateCategory)))
	mux.Handle("/api/categories/delete", middleware.AuthMiddleware(http.HandlerFunc(s.handleDeleteCategory)))

	return withCORS(mux)
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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

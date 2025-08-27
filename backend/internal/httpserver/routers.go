package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google-auth-demo/backend/internal/service"

	"github.com/joho/godotenv"
)

var (
	port         string
	redirectBase string
	frontendURL  string
)

type Server struct {
	router *http.ServeMux
	svc    *service.Service
}

func New(svc *service.Service) *Server {
	_ = godotenv.Load()

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redirectBase = os.Getenv("GOOGLE_REDIRECT_URI_BASE")
	if redirectBase == "" {
		redirectBase = "http://localhost"
	}

	frontendURL = os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	server := &Server{
		router: http.NewServeMux(),
		svc:    svc,
	}

	server.routes()
	return server
}

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleHome)
	s.router.HandleFunc("/login", s.handleLogin)
	s.router.HandleFunc("/oauth2callback", s.handleCallback)
}

func (s *Server) Start() error {
	addr := ":" + port
	fmt.Printf("🌐 Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(addr, s.router))
	return nil
}

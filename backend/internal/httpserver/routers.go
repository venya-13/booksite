package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	port            string
	clientID        string
	clientSecret    string
	redirectURIBase string
	frontendURL     string
)

type Service interface {
	HandleCallback() error
}

type Server struct {
	// port int
	router *http.ServeMux
	svc    Service
}

func New(svc Service) *Server {

	_ = godotenv.Load() // не падаем, если .env нет — можно брать из окружения

	port = os.Getenv("PORT")
	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURIBase = os.Getenv("GOOGLE_REDIRECT_URI_BASE")
	if redirectURIBase == "" {
		redirectURIBase = "http://localhost"
	}
	frontendURL = os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	router := http.NewServeMux()
	server := &Server{
		router: router,
		svc:    svc,
	}

	server.routes()

	return server
}

func (s *Server) routes() {
	s.router.HandleFunc("/", handleHome)
	s.router.HandleFunc("/login", handleLogin)
	s.router.HandleFunc("/oauth2callback", HandleCallback)
}

func (s *Server) Start() error {

	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	fmt.Println("🌐 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(addr, s.router))
	return nil
}

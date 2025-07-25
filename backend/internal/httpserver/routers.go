package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI = os.Getenv("GOOGLE_REDIRECT_URI")

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
	fmt.Println("🌐 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", s.router))
	return nil
}

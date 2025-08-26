package httpserver

import (
	"fmt"
	"net/http"
)

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Google Auth Demo!")
	fmt.Fprintln(w, "Click here to <a href='/login'>Login with Google</a>")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := s.svc.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code found in callback", http.StatusBadRequest)
		return
	}

	userJson, err := s.svc.HandleCallback(code)
	if err != nil {
		http.Error(w, "Callback error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := s.svc.GetFrontendURL(userJson)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

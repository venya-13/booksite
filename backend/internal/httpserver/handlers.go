package httpserver

import (
	"encoding/json"
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

func (s *Server) handleSomeProtectedAPI(w http.ResponseWriter, r *http.Request) {
	googleID := "some-google-id" // In real scenario, extract from JWT or session
	// need jwt to funtion properly

	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "Google API error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User profile: %+v", profile)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	googleID := r.URL.Query().Get("google_id")
	if googleID == "" {
		http.Error(w, "google_id required", http.StatusBadRequest)
		return
	}

	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "failed to refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "failed to fetch profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo["id"] = googleID
	userInfo["access_token"] = accessToken
	if err := s.svc.SaveUser(userInfo); err != nil {
		http.Error(w, "failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token": accessToken,
		"user":         userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) GetProfile(w http.ResponseWriter, r *http.Request) {
	googleID := r.URL.Query().Get("google_id")
	if googleID == "" {
		http.Error(w, "google_id required", http.StatusBadRequest)
		return
	}

	// gет access token
	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "failed to ensure token", http.StatusInternalServerError)
		return
	}

	// get profile
	profile, err := s.svc.FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "failed to fetch profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(profile)
}

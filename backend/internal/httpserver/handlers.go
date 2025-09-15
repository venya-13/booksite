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
	googleID := "какой-то_google_id_из_сессии_или_jwt"

	accessToken, err := s.svc.EnsureAccessToken(googleID)
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), 500)
		return
	}

	// we can now use this access token to call Google APIs
	profile, err := s.svc.OAuth().FetchProfile(accessToken)
	if err != nil {
		http.Error(w, "Google API error: "+err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "User profile: %+v", profile)
}

func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		http.Error(w, "refresh_token required", http.StatusBadRequest)
		return
	}

	tokenData, err := s.svc.OAuth().RefreshAccessToken(refreshToken)
	if err != nil {
		http.Error(w, "failed to refresh token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := s.svc.OAuth().FetchProfile(tokenData.AccessToken)
	if err != nil {
		http.Error(w, "failed to fetch profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// save new user info and rotated refresh token
	userInfo["refresh_token"] = tokenData.RefreshToken
	if err := s.svc.SaveUser(userInfo); err != nil {
		http.Error(w, "failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// give new tokens and user info back to client
	response := map[string]interface{}{
		"access_token":  tokenData.AccessToken,
		"refresh_token": tokenData.RefreshToken,
		"id_token":      tokenData.IdToken,
		"expires_in":    tokenData.ExpiresIn,
		"user":          userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

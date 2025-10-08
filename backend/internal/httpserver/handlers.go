package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
		slog.Warn("Missing code in callback", slog.String("url", r.URL.String()))
		http.Error(w, "No code found in callback", http.StatusBadRequest)
		return
	}

	slog.Info("OAuth callback received", slog.String("code", code[:6]+"...")) // log only first 6 chars for privacy

	resp, err := s.svc.HandleCallback(code)
	if err != nil {
		slog.Error("HandleCallback failed", slog.String("error", err.Error()))
		http.Error(w, "Callback error", http.StatusInternalServerError)
		return
	}

	id, _ := resp.User["id"].(string)
	email, _ := resp.User["email"].(string)
	slog.Info("User authenticated successfully",
		slog.String("google_id", id),
		slog.String("email", email),
	)

	redirectURL := s.svc.GetFrontendURL(resp.JWT)
	slog.Info("Redirecting user to frontend", slog.String("url", redirectURL))

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (s *Server) handleGoogleProfile(w http.ResponseWriter, r *http.Request) {
	googleID, _ := r.Context().Value("user_id").(string)
	if googleID == "" {
		http.Error(w, "missing user_id in context", http.StatusUnauthorized)
		return
	}

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

	// return JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
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

func (s *Server) handleProtected(w http.ResponseWriter, r *http.Request) {
	// values come from middleware
	userID, _ := r.Context().Value("user_id").(string)
	email, _ := r.Context().Value("email").(string)
	isAdmin, _ := r.Context().Value("is_admin").(bool)

	resp := map[string]interface{}{
		"message":  "This is a protected route",
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleJWTRefresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "refresh_token missing", http.StatusUnauthorized)
		return
	}

	resp, err := s.svc.RefreshJWT(cookie.Value)
	if err != nil {
		http.Error(w, "invalid refresh token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// set new tokens in secure cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.svc.JWTTTL.Seconds()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(s.svc.RefreshTTL.Seconds()),
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "JWT refreshed successfully"}`))
}

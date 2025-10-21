package service

import (
	"fmt"
	"google-auth-demo/backend/internal/jwt"
	"log/slog"
	"net/url"
	"time"
)

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type OAuth interface {
	GetAuthURL() string
	ExchangeCode(string) (*TokenData, error)
	RefreshAccessToken(string) (*TokenData, error)
	FetchProfile(string) (map[string]interface{}, error)
}

type Repository interface {
	SaveOrUpdate(user map[string]interface{}) error
	GetUserByGoogleID(googleID string) (map[string]interface{}, error)
}

type (
	Service struct {
		OAuth       OAuth
		Repo        Repository
		FrontendURL string
		JWTTTL      time.Duration
		RefreshTTL  time.Duration
	}

	Config struct {
		FrontendURL string `env:"FRONTEND_URL"`
	}
)

type AuthResponse struct {
	User         map[string]interface{} `json:"user"`
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	JWT          string                 `json:"jwt"`
}

func New(frontendURL string, oauth OAuth, repo Repository, jwtTTL, refreshTTL time.Duration) *Service {
	return &Service{
		FrontendURL: frontendURL,
		OAuth:       oauth,
		Repo:        repo,
		JWTTTL:      jwtTTL,
		RefreshTTL:  refreshTTL,
	}
}

func (s *Service) GetAuthURL() string {
	return s.OAuth.GetAuthURL()
}

func (s *Service) HandleCallback(code string) (*AuthResponse, error) {
	slog.Info("Handling OAuth callback")

	tokenData, err := s.OAuth.ExchangeCode(code)
	if err != nil {
		slog.Error("Failed to exchange code", slog.String("error", err.Error()))
		return nil, err
	}

	userInfo, err := s.OAuth.FetchProfile(tokenData.AccessToken)
	if err != nil {
		slog.Error("Failed to fetch Google profile", slog.String("error", err.Error()))
		return nil, err
	}

	id, _ := userInfo["id"].(string)
	email, _ := userInfo["email"].(string)
	slog.Info("Fetched Google profile", slog.String("google_id", id), slog.String("email", email))

	userInfo["access_token"] = tokenData.AccessToken
	userInfo["refresh_token"] = tokenData.RefreshToken
	userInfo["token_expiry"] = time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)

	existingUser, err := s.Repo.GetUserByGoogleID(id)
	if err != nil {
		slog.Warn("User not found in DB, will be created", slog.String("google_id", id))
	}

	if err := s.Repo.SaveOrUpdate(userInfo); err != nil {
		slog.Error("DB save/update failed", slog.String("google_id", id), slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("User saved to database", slog.String("google_id", id))

	// isAdmin or not
	isAdmin := false
	if existingUser != nil {
		if v, ok := existingUser["is_admin"].(bool); ok {
			isAdmin = v
		}
	}

	jwtToken, _, err := s.GenerateTokens(id, email, isAdmin)
	if err != nil {
		slog.Error("Failed to generate JWT", slog.String("google_id", id), slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("JWT token generated", slog.String("google_id", id), slog.Bool("is_admin", isAdmin))

	return &AuthResponse{
		User:         userInfo,
		AccessToken:  tokenData.AccessToken,
		RefreshToken: userInfo["refresh_token"].(string),
		JWT:          jwtToken,
	}, nil
}

func (s *Service) GetFrontendURL(jwtToken string) string {
	return s.FrontendURL + "?token=" + url.QueryEscape(jwtToken)
}

func (s *Service) EnsureAccessToken(googleID string) (string, error) {
	user, err := s.Repo.GetUserByGoogleID(googleID)
	if err != nil {
		return "", err
	}

	expiryInterface, ok := user["token_expiry"]
	if !ok {
		return "", fmt.Errorf("token expiry missing for user %s", googleID)
	}

	expiry, ok := expiryInterface.(time.Time)
	if !ok {
		return "", fmt.Errorf("invalid token expiry type for user %s", googleID)
	}

	accessToken, ok := user["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("invalid access token type for user %s", googleID)
	}

	// if token is still valid, return it
	if time.Now().Before(expiry) {
		return accessToken, nil
	}

	// update token
	refreshToken, ok := user["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return "", fmt.Errorf("no refresh token available for user %s", googleID)
	}

	newToken, err := s.OAuth.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", err
	}

	// if google did not return a new refresh token, keep the old one
	if newToken.RefreshToken == "" {
		newToken.RefreshToken = refreshToken
	}

	// updaate user info
	user["access_token"] = newToken.AccessToken
	user["refresh_token"] = newToken.RefreshToken
	user["token_expiry"] = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)

	if err := s.Repo.SaveOrUpdate(user); err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

func (s *Service) FetchProfile(accessToken string) (map[string]interface{}, error) {
	return s.OAuth.FetchProfile(accessToken)
}

func (s *Service) SaveUser(user map[string]interface{}) error {
	return s.Repo.SaveOrUpdate(user)
}

func (s *Service) GenerateTokens(googleID, email string, isAdmin bool) (string, string, error) {
	accessToken, err := jwt.GenerateToken(googleID, email, isAdmin, s.JWTTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken(googleID, s.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshJWT(refreshToken string) (*AuthResponse, error) {
	claims, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	newAccessToken, newRefreshToken, err := s.GenerateTokens(claims.GoogleID, claims.Email, claims.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

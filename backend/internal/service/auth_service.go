package service

import (
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

type AuthResponse struct {
	User         map[string]interface{} `json:"user"`
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	JWT          string                 `json:"jwt"`
}

type Service struct {
	OAuth       OAuth
	Repo        Repository
	FrontendURL string
	JWTTTL      time.Duration
	RefreshTTL  time.Duration
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
	userInfo["google_id"] = id
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

	jwtToken, refreshToken, err := s.GenerateTokens(id, email, isAdmin)

	if err != nil {
		slog.Error("Failed to generate JWT", slog.String("google_id", id), slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("JWT token generated", slog.String("google_id", id), slog.Bool("is_admin", isAdmin))

	return &AuthResponse{
		User:         userInfo,
		AccessToken:  tokenData.AccessToken,
		RefreshToken: refreshToken,
		JWT:          jwtToken,
	}, nil
}

func (s *Service) GetFrontendURL(jwtToken string) string {
	return s.FrontendURL + "?token=" + url.QueryEscape(jwtToken)
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

func (s *Service) FetchProfile(accessToken string) (map[string]interface{}, error) {
	return s.OAuth.FetchProfile(accessToken)
}

func (s *Service) SaveUser(user map[string]interface{}) error {
	return s.Repo.SaveOrUpdate(user)
}

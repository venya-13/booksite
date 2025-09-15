package service

import (
	"encoding/json"
	"google-auth-demo/backend/internal/repo"
	"net/url"
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
}

type (
	Service struct {
		oauth       OAuth
		repo        Repository
		frontendURL string
	}

	Config struct {
		FrontendURL string `env:"FRONTEND_URL"`
	}
)

func New(config Config, oauth OAuth, repo Repository) *Service {
	return &Service{
		frontendURL: config.FrontendURL,
		oauth:       oauth,
		repo:        repo,
	}
}

func (s *Service) GetAuthURL() string {
	return s.oauth.GetAuthURL()
}

func (s *Service) HandleCallback(code string) (string, error) {
	tokenData, err := s.oauth.ExchangeCode(code)
	if err != nil {
		return "", err
	}

	userInfo, err := s.oauth.FetchProfile(tokenData.AccessToken)
	if err != nil {
		return "", err
	}

	// add refresh token to user info for saving in DB
	userInfo["refresh_token"] = tokenData.RefreshToken

	if err := s.repo.SaveOrUpdate(userInfo); err != nil {
		return "", err
	}

	userJson, _ := json.Marshal(userInfo)
	return string(userJson), nil
}

func (s *Service) GetFrontendURL(userJson string) string {
	return s.frontendURL + "?user=" + url.QueryEscape(userJson)
}

func (s *Service) EnsureAccessToken(googleID string) (string, error) {
	// 1. get refresh_token from DB
	refreshToken, err := s.repo.(*repo.PostgresRepo).GetRefreshTokenByGoogleID(googleID)
	if err != nil {
		return "", err
	}

	// request new access_token using refresh_token
	newToken, err := s.oauth.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", err
	}

	// save new refresh_token if it was rotated
	userInfo := map[string]interface{}{
		"id":            googleID,
		"refresh_token": newToken.RefreshToken,
	}
	_ = s.repo.SaveOrUpdate(userInfo)

	return newToken.AccessToken, nil
}

func (s *Service) OAuth() OAuth {
	return s.oauth
}

func (s *Service) SaveUser(user map[string]interface{}) error {
	return s.repo.SaveOrUpdate(user)
}

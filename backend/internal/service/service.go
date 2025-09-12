package service

import (
	"encoding/json"
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
		return "error:", err
	}

	userInfo, err := s.oauth.FetchProfile(tokenData.AccessToken)
	if err != nil {
		return "", err
	}

	if err := s.repo.SaveOrUpdate(userInfo); err != nil {
		return "", err
	}

	userJson, _ := json.Marshal(userInfo)
	return string(userJson), nil
}

func (s *Service) GetFrontendURL(userJson string) string {
	return s.frontendURL + "?user=" + url.QueryEscape(userJson)
}

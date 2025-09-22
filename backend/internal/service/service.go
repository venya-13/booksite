package service

import (
	"encoding/json"
	"fmt"
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

	// get already exist
	existingUser, _ := s.repo.GetUserByGoogleID(userInfo["id"].(string))

	// update token info
	userInfo["access_token"] = tokenData.AccessToken
	userInfo["token_expiry"] = time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)

	if tokenData.RefreshToken != "" {
		userInfo["refresh_token"] = tokenData.RefreshToken
	} else if existingUser != nil {
		userInfo["refresh_token"] = existingUser["refresh_token"]
	}

	id, ok := userInfo["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("user id missing in profile")
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		email = ""
	}

	name, ok := userInfo["name"].(string)
	if !ok {
		name = ""
	}

	picture, ok := userInfo["picture"].(string)
	if !ok {
		picture = ""
	}

	userInfo["id"] = id
	userInfo["email"] = email
	userInfo["name"] = name
	userInfo["picture"] = picture

	// save to database
	if err := s.repo.SaveOrUpdate(userInfo); err != nil {
		return "", err
	}

	// back to frontend
	userJson, _ := json.Marshal(userInfo)
	return string(userJson), nil
}

func (s *Service) GetFrontendURL(userJson string) string {
	return s.frontendURL + "?user=" + url.QueryEscape(userJson)
}

func (s *Service) EnsureAccessToken(googleID string) (string, error) {
	user, err := s.repo.GetUserByGoogleID(googleID)
	if err != nil {
		return "", err
	}

	expiry := user["token_expiry"].(time.Time)
	accessToken := user["access_token"].(string)

	// if token is still valid, return it
	if time.Now().Before(expiry) {
		return accessToken, nil
	}

	// update token
	refreshToken, ok := user["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return "", fmt.Errorf("no refresh token available for user %s", googleID)
	}

	newToken, err := s.oauth.RefreshAccessToken(refreshToken)
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

	if err := s.repo.SaveOrUpdate(user); err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

func (s *Service) FetchProfile(accessToken string) (map[string]interface{}, error) {
	return s.oauth.FetchProfile(accessToken)
}

func (s *Service) SaveUser(user map[string]interface{}) error {
	return s.repo.SaveOrUpdate(user)
}

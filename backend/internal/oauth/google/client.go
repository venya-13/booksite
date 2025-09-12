package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"google-auth-demo/backend/internal/service"
)

type (
	Client struct {
		redirectURL  string
		clientID     string
		clientSecret string

		// Temporary in-memory token storage (later will be replaced with DB)
		tokenData *service.TokenData
		expiry    time.Time
	}

	Config struct {
		RedirectBaseURL string `env:"HTTP_SERVER_REDIRECT_BASE_URL"`
		ClientID        string `env:"GOOGLE_CLIENT_ID"`
		ClientSecret    string `env:"GOOGLE_CLIENT_SECRET"`
	}
)

func New(config Config) *Client {
	redirectURL := fmt.Sprintf("%s/oauth2callback", config.RedirectBaseURL)

	return &Client{
		redirectURL:  redirectURL,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
	}
}

// GetAuthURL builds Google OAuth2 authorization URL
func (c *Client) GetAuthURL() string {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
		"client_id":     {c.clientID},
		"redirect_uri":  {c.redirectURL},
		"response_type": {"code"},
		"scope":         {"email profile"},
		"state":         {"random123"},
		"access_type":   {"offline"},
		"prompt":        {"select_account"},
	}.Encode()

	return authURL
}

// ExchangeCode exchanges authorization code for access and refresh tokens
func (c *Client) ExchangeCode(code string) (*service.TokenData, error) {
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"redirect_uri":  {c.redirectURL},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenData service.TokenData
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	// store token in memory
	c.tokenData = &tokenData
	c.expiry = time.Now().Add(time.Second * time.Duration(tokenData.ExpiresIn))

	return &tokenData, nil
}

// RefreshAccessToken requests a new access token using the refresh token
func (c *Client) RefreshAccessToken(refreshToken string) (*service.TokenData, error) {
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenData service.TokenData
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	// update in-memory token storage
	if tokenData.RefreshToken == "" && c.tokenData != nil {
		tokenData.RefreshToken = c.tokenData.RefreshToken
	}

	c.tokenData = &tokenData
	c.expiry = time.Now().Add(time.Second * time.Duration(tokenData.ExpiresIn))

	return &tokenData, nil
}

func (c *Client) EnsureAccessToken() (string, error) {
	if c.tokenData == nil {
		return "", fmt.Errorf("token is not initialized")
	}

	// check expiry time
	if time.Now().After(c.expiry) {
		newToken, err := c.RefreshAccessToken(c.tokenData.RefreshToken)
		if err != nil {
			return "", err
		}
		return newToken.AccessToken, nil
	}

	return c.tokenData.AccessToken, nil
}

func (c *Client) FetchProfile(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

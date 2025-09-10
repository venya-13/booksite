package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"google-auth-demo/backend/internal/service"
)

type (
	Client struct {
		redirectURL  string
		clientID     string
		clientSecret string
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

	return &tokenData, nil
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

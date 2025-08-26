package oauth

import (
	"encoding/json"
	"fmt"
	"google-auth-demo/backend/internal/service"
	"net/http"
	"net/url"
	"os"
)

type GoogleOAuth struct{}

func NewGoogleOAuth() *GoogleOAuth {
	return &GoogleOAuth{}
}

func (g *GoogleOAuth) GetAuthURL() string {
	redirectURI := fmt.Sprintf("%s:%s/oauth2callback",
		os.Getenv("GOOGLE_REDIRECT_URI_BASE"),
		os.Getenv("PORT"))

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"redirect_uri":  {redirectURI},
		"response_type": {"code"},
		"scope":         {"email profile"},
		"state":         {"random123"},
		"access_type":   {"offline"},
		"prompt":        {"select_account"},
	}.Encode()

	return authURL
}

func (g *GoogleOAuth) ExchangeCode(code string) (*service.TokenData, error) {
	redirectURI := fmt.Sprintf("%s:%s/oauth2callback",
		os.Getenv("GOOGLE_REDIRECT_URI_BASE"),
		os.Getenv("PORT"))

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"redirect_uri":  {redirectURI},
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

func (g *GoogleOAuth) FetchProfile(accessToken string) (map[string]interface{}, error) {
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

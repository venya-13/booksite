package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var (
	clientID     string
	clientSecret string
	redirectURI  string
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI = os.Getenv("GOOGLE_REDIRECT_URI")

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/oauth2callback", handleCallback)

	fmt.Println("🌐 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Google Auth Demo!")
	fmt.Fprintln(w, "Click here to <a href='/login'>Login with Google</a>")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {redirectURI},
		"response_type": {"code"},
		"scope":         {"email profile"},  // not obligatory
		"state":         {"random123"},      // For CSRF protection (optional)
		"access_type":   {"offline"},        // Request offline access to get a refresh token
		"prompt":        {"select_account"}, // give you the option to select an account if multiple accounts are logged in
	}.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code found in callback", http.StatusBadRequest)
		return
	}

	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})

	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	json.NewDecoder(tokenResp.Body).Decode(&tokenData)

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	client := &http.Client{}
	userInfoResp, err := client.Do(req)

	var userInfo map[string]interface{}
	json.NewDecoder(userInfoResp.Body).Decode(&userInfo)

	userJson, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, "Failed to encode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to React frontend with user info in query string
	redirectURL := "http://localhost:5173?user=" + url.QueryEscape(string(userJson))
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

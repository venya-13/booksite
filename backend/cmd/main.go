package main

import (
	"fmt"
	"google-auth-demo/backend/internal/httpserver"
	"google-auth-demo/backend/internal/oauth"
	"google-auth-demo/backend/internal/repo"
	"google-auth-demo/backend/internal/service"
)

func main() {
	fmt.Println("Start the program")

	oauthGoogle := oauth.NewGoogleOAuth()
	repository := repo.NewMockRepo()
	svc := service.New(oauthGoogle, repository)

	server := httpserver.New(svc)
	server.Start()
}

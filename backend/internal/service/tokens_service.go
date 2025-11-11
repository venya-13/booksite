package service

import (
	"fmt"
	"google-auth-demo/backend/internal/jwt"
	"time"
)

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

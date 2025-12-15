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

	// still valid
	if time.Now().Before(expiry) {
		return accessToken, nil
	}

	// refresh from Google
	refreshToken, ok := user["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return "", fmt.Errorf("no refresh token available for user %s", googleID)
	}

	newToken, err := s.OAuth.RefreshAccessToken(refreshToken)
	if err != nil {
		return "", err
	}

	if newToken.RefreshToken == "" {
		newToken.RefreshToken = refreshToken
	}

	user["access_token"] = newToken.AccessToken
	user["refresh_token"] = newToken.RefreshToken
	user["token_expiry"] = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)

	if err := s.Repo.SaveOrUpdate(user); err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}

func (s *Service) RefreshJWT(refreshToken string) (*AuthResponse, error) {
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.Repo.GetUserByGoogleID(claims.GoogleID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	email, _ := user["email"].(string)
	isAdmin, _ := user["is_admin"].(bool)

	newAccessToken, newRefreshToken, err := s.GenerateTokens(
		claims.GoogleID,
		email,
		isAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	user["access_token"] = newAccessToken
	user["refresh_token"] = newRefreshToken
	user["token_expiry"] = time.Now().Add(s.JWTTTL)

	if err := s.Repo.SaveOrUpdate(user); err != nil {
		return nil, fmt.Errorf("failed to save updated tokens: %w", err)
	}

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

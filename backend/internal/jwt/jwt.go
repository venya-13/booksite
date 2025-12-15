package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret []byte
var refreshSecret []byte

var (
	ErrNoAccessSecret  = errors.New("access secret is not set")
	ErrNoRefreshSecret = errors.New("refresh secret is not set")
)

type Claims struct {
	GoogleID string `json:"google_id"`
	Email    string `json:"email,omitempty"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func SetSecrets(access, refresh string) {
	accessSecret = []byte(access)
	refreshSecret = []byte(refresh)
}

func mustHaveAccessSecret() error {
	if len(accessSecret) == 0 {
		return ErrNoAccessSecret
	}
	return nil
}

func mustHaveRefreshSecret() error {
	if len(refreshSecret) == 0 {
		return ErrNoRefreshSecret
	}
	return nil
}

// ACCESS TOKEN
func GenerateAccessToken(googleID, email string, isAdmin bool, duration time.Duration) (string, error) {
	if err := mustHaveAccessSecret(); err != nil {
		return "", err
	}

	claims := &Claims{
		GoogleID: googleID,
		Email:    email,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// REFRESH TOKEN
func GenerateRefreshToken(googleID string, duration time.Duration) (string, error) {
	if err := mustHaveRefreshSecret(); err != nil {
		return "", err
	}

	claims := &Claims{
		GoogleID: googleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// VALIDATION
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	if err := mustHaveAccessSecret(); err != nil {
		return nil, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrTokenUnverifiable
		}
		return accessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	if err := mustHaveRefreshSecret(); err != nil {
		return nil, err
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrTokenUnverifiable
		}
		return refreshSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenUnverifiable
	}

	return claims, nil
}

func GenerateToken(googleID, email string, isAdmin bool, duration time.Duration) (string, error) {
	return GenerateAccessToken(googleID, email, isAdmin, duration)
}

func ValidateToken(token string) (*Claims, error) {
	return ValidateAccessToken(token)
}

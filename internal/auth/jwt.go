package auth

import (
	"errors"
	"time"

	"github.com/0cd/go-ecom/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(env.GetString("ACCESS_SECRET", "your-access-secret"))
	refreshSecret = []byte(env.GetString("REFRESH_SECRET", "your-refresh-secret"))
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int64) (string, error) {
	expirationDate := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID int64) (string, error) {
	expirationDate := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ParseAccessToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, accessSecret)
}

func ParseRefreshToken(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr, refreshSecret)
}

func parseToken(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

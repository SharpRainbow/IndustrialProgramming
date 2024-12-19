package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"onlineShop/internal/models"
	"strings"
	"time"
)

var jwtKey = []byte("my_secret_key")

func GenerateToken(username string, role *models.Role) (string, error) {
	if role == nil {
		return "", errors.New("user role not found")
	}
	iat := time.Now()
	claims := jwt.StandardClaims{
		Audience:  strings.TrimSpace(role.Name),
		Subject:   username,
		IssuedAt:  iat.Unix(),
		ExpiresAt: iat.Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GetKey() []byte {
	return jwtKey
}

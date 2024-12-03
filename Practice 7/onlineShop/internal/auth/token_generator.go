package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("my_secret_key")

func GenerateToken(username string) (string, error) {
	iat := time.Now()
	claims := jwt.StandardClaims{
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

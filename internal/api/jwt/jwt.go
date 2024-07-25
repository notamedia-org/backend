package tokenHandler

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/notamedia-org/backend/internal/config"
)

func CreateToken(env *config.Config) (string, error) {
	signingKey := []byte(env.JwtSecret)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}

func VerifyToken(env *config.Config, tokenString string) (bool, error) {
	claims := &jwt.RegisteredClaims{}

	signingKey := []byte(env.JwtSecret)
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return false, nil
	}

	return true, err
}

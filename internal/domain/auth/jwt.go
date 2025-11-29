package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret  string
	expires time.Duration
}

func NewJWTManager(secret string, expires time.Duration) *JWTManager {
	return &JWTManager{secret: secret, expires: expires}
}

func (j *JWTManager) Sign(id int, role UserRole) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":   id,
			"role": string(role),
			"exp":  time.Now().Add(j.expires).Unix(),
			"iat":  time.Now().Unix(),
		},
	)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) Verify(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

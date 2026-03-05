package middleware

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type TokenService struct {
	auth *jwtauth.JWTAuth
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		auth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func (t *TokenService) Generate(sessionId string) (string, error) {

	_, tokenString, err := t.auth.Encode(map[string]interface{}{
		"session_id": sessionId,
		"exp":        time.Now().Add(24 * time.Hour),
	})

	return tokenString, err
}

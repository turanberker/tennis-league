package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/turanberker/tennis-league-service/internal/domain/session"
)

type TokenService struct {
	auth              *jwtauth.JWTAuth
	sessionRepository session.Repository
}

func NewTokenService(secret string, sessionRepository session.Repository) *TokenService {
	return &TokenService{
		auth:              jwtauth.New("HS256", []byte(secret), nil),
		sessionRepository: sessionRepository,
	}
}

func (t *TokenService) GenerateAccessTokenAndSetCookie(c *gin.Context, sessionId string) (string, error) {

	_, tokenString, err := t.auth.Encode(map[string]interface{}{
		"session_id": sessionId,
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	})
	// CSRF Koruması için SameSite modunu ayarla
	c.SetSameSite(http.SameSiteLaxMode)
	// 2. Access Token Cookie (Her istekte gönderilmeli)
	c.SetCookie(
		"access_token",
		tokenString,
		3600,  // 1 saat (Access Token süresiyle uyumlu)
		"/",   // Tüm endpoint'lerde geçerli
		"",    // Domain
		false, // Prod'da true (HTTPS)
		true,  // HttpOnly: JS erişemez
	)
	return tokenString, err
}

func (t *TokenService) GenerateRefreshTokenAndSetCookie(c *gin.Context, sessionId string) (string, error) {
	// Refresh Token genellikle daha uzun ömürlüdür (örn. 7 gün)
	_, tokenString, err := t.auth.Encode(map[string]interface{}{
		"session_id": sessionId,
		"type":       "refresh", // Token tipini belirtmek güvenliği artırır
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token", // isim
		tokenString,     // değer
		3600*24*7,       // 7 gün (saniye cinsinden)
		"/auth/refresh", // ÖNEMLİ: Sadece refresh endpointine gönderilsin
		"",              // domain
		false,           // local'de çalıştığın için false, prod'da TRUE olmalı (HTTPS)
		true,            // HTTP_ONLY: JavaScript erişemez (XSS koruması)
	)

	return tokenString, err
}

func (t *TokenService) ValidateAndRefreshAndSetAccessCookie(c *gin.Context, refreshToken string) (string, error) {
	// 1. Refresh Token'ı decode et
	token, err := t.auth.Decode(refreshToken)
	if err != nil {
		return "", errors.New("refresh token geçersiz")
	}

	// 2. session_id'yi al
	claims := token.PrivateClaims()
	sId, ok := claims["session_id"].(string)
	if !ok {
		return "", errors.New("token içinde session_id bulunamadı")
	}

	// 3. Redis'te bu session hala var mı? (Get metodunu kullanıyoruz)
	// Access Token süresi (1 saat) dolsa bile, Redis TTL 7 gün olduğu için
	// kullanıcı logout yapmadığı sürece bu session burada duracaktır.
	sess, err := t.sessionRepository.Get(c.Request.Context(), sId)
	if err != nil || sess == nil {
		return "", errors.New("oturum süresi dolmuş, lütfen tekrar giriş yapın")
	}

	// 4. (Opsiyonel) Her refresh işleminde Redis süresini tekrar uzatabilirsin
	_ = t.sessionRepository.Refresh(c.Request.Context(), sId)

	// 5. Yeni kısa ömürlü Access Token üret (Örn: 1 saatlik)
	return t.GenerateAccessTokenAndSetCookie(c, sId)
}

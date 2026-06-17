package authmiddleware

import (
	"errors"
	"net/http"
	"strings"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/security/dto"
	"tennis-league/common/security/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
)

var errSessionExpired = "AUTH_102"
var insufficient_permissions = "AUTH_100"

type AuthMiddleware struct {
	tokenAuth         *jwtauth.JWTAuth
	sessionRepository repository.SessionGetterRepository
}

func NewAuthMiddleware(secret string, sessionRepository repository.SessionGetterRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenAuth:         jwtauth.New("HS256", []byte(secret), nil),
		sessionRepository: sessionRepository,
	}
}

func (a *AuthMiddleware) GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Önce Cookie'den access_token'ı oku
		tokenString, err := c.Cookie("access_token")

		if err != nil {
			c.Next() // Sessiz devam, RequireAuth yakalayacak
			return
		}

		// Token doğrulaması
		token, err := jwtauth.VerifyToken(a.tokenAuth, tokenString)
		if err != nil {
			// Hata olsa bile sessiz kalıyoruz, sadece devam ediyoruz
			c.Next()
			return
		}

		claims := token.PrivateClaims()
		sessionId, ok := claims["session_id"].(string)
		if !ok || sessionId == "" {
			c.Next()
			return
		}

		// Session kontrolü
		session, err := a.sessionRepository.Get(c, sessionId)
		if err != nil || session == nil {
			c.Next()
			return
		}

		// Buraya geldiyse her şey yolunda demektir, context'i doldurabiliriz
		c.Set("session_id", sessionId)
		c.Set("Role", dto.Role(session.Role))
		c.Set("UserId", session.UserId)

		if session.PlayerId != nil {
			c.Set("PlayerId", *session.PlayerId)
		} else {
			c.Set("PlayerId", "")
		}

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GetToken bir session_id bulup set etti mi?
		sessionId, _ := c.Get("session_id")

		if sessionId == "" || sessionId == nil {
			err := &customerror.BusinnesException{
				StatusCode: http.StatusUnauthorized,
				ErrorCode:  errSessionExpired, // Burada AUTH_102 dönersen React atar
				Message:    "Oturumunuzun süresi dolmuş veya geçersiz. Lütfen tekrar giriş yapın.",
			}
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(roles ...dto.Role) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, hasRole := c.Get("Role")

		roleStrings := make([]string, len(roles))
		for i, r := range roles {
			roleStrings[i] = string(r)
		}

		if !hasRole {

			// Özel bir business hatası oluşturuyoruz
			err := &customerror.BusinnesException{
				StatusCode: http.StatusForbidden,
				ErrorCode:  insufficient_permissions,
				Message:    "Kullanıcı şu rollerden birine sahip olmalı: " + strings.Join(roleStrings, ", "),
			}
			c.Error(err) // Hatayı Gin'in listesine ekle
			c.Abort()    // İsteği durdur (Handler'a gitmesin)
			return
		}

		userRole, ok := roleValue.(dto.Role)
		if !ok {
			err := customerror.NewInternalError(errors.New("Geçersiz Rol Tipi"))
			c.Error(err)
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}
		err := &customerror.BusinnesException{
			StatusCode: http.StatusForbidden,
			ErrorCode:  insufficient_permissions,
			Message:    "Kullanıcı şu rollerden birine sahip olmalı: " + strings.Join(roleStrings, ", "),
		}
		c.Error(err) // Hatayı Gin'in listesine ekle
		c.Abort()    // İsteği durdur (Handler'a gitmesin)
	}
}

func GetUserIdFromContext(c *gin.Context) (string, bool) {
	userIdValue, exists := c.Get("UserId")
	return userIdValue.(string), exists
}

func GetPlayerIdFromContext(c *gin.Context) (string, bool) {
	playerIdIdValue, exists := c.Get("PlayerId")
	return playerIdIdValue.(string), exists
}

package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/redis"
)

type AuthMiddleware struct {
	tokenAuth         *jwtauth.JWTAuth
	sessionRepository *redis.SessionRepository
}

func NewAuthMiddleware(secret string, sessionRepository *redis.SessionRepository) *AuthMiddleware {
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
		c.Set("Role", user.Role(session.Role))
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
				ErrorCode:  customerror.ErrSessionExpired, // Burada AUTH_102 dönersen React atar
				Message:    "Oturumunuzun süresi dolmuş veya geçersiz. Lütfen tekrar giriş yapın.",
			}
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRole(roles ...user.Role) gin.HandlerFunc {
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
				ErrorCode:  customerror.INSUFFICIENT_PERMISSIONS,
				Message:    "Kullanıcı şu rollerden birine sahip olmalı: " + strings.Join(roleStrings, ", "),
			}
			c.Error(err) // Hatayı Gin'in listesine ekle
			c.Abort()    // İsteği durdur (Handler'a gitmesin)
			return
		}

		userRole, ok := roleValue.(user.Role)
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
			ErrorCode:  customerror.INSUFFICIENT_PERMISSIONS,
			Message:    "Kullanıcı şu rollerden birine sahip olmalı: " + strings.Join(roleStrings, ", "),
		}
		c.Error(err) // Hatayı Gin'in listesine ekle
		c.Abort()    // İsteği durdur (Handler'a gitmesin)
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Önce handler çalışsın

		// Eğer bir hata varsa
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 1. Business Exception Kontrolü
			var businessErr *customerror.BusinnesException
			if errors.As(err, &businessErr) {
				c.JSON(businessErr.StatusCode, delivery.NewBusinnesErrorResponse(*businessErr))
				return // Yanıt gönderildikten sonra durmalı
			}

			// 2. Internal Exception Kontrolü
			var internalErr *customerror.InternalException
			if errors.As(err, &internalErr) {
				// appErr.Message kontrolü yerine direkt RawError kontrolü daha mantıklı olabilir
				if internalErr.RawError != nil {
					log.Printf("[SİSTEM HATASI]: %v", internalErr.RawError)
				}

				// Kullanıcıya teknik detay değil, InternalException içindeki güvenli mesajı dönüyoruz
				c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
				return
			}

			// 3. Beklenmedik Hatalar
			log.Printf("[KRİTİK HATA]: %v", err)
			c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
		}
	}
}

func GetUserIdFromContext(c *gin.Context) (string, bool) {
	userIdValue, exists := c.Get("UserId")
	return userIdValue.(string), exists
}
func AddCacheControlHeader(maxAge int, cacheType Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dinamik olarak header değerini oluşturuyoruz
		headerValue := fmt.Sprintf("%s, max-age=%d", cacheType, maxAge)

		c.Header("Cache-Control", headerValue)

		// Bir sonraki handler'a geçişi sağlar
		c.Next()
	}
}

type Type string

const (
	TYPE_PUBLIC  Type = "public"
	TYPE_PRIVATE Type = "private"
)

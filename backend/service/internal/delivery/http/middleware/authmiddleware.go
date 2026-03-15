package middleware

import (
	"errors"
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

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
			return // 🔥 BURASI EKSİKTİ
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtauth.VerifyToken(a.tokenAuth, tokenString)
		if err != nil {
			err := &customerror.BusinnesException{
				StatusCode: http.StatusUnauthorized,
				ErrorCode:  customerror.ErrSessionExpired,
				Message:    "Oturumunuz bitmiştir.",
			}
			c.Error(err) // Hatayı Gin'in listesine ekle
			c.Abort()    // İsteği durdur (Handler'a gitmesin)
			return
		}

		claims := token.PrivateClaims()
		sessionId := claims["session_id"].(string)
		// Context'e user bilgilerini koy
		c.Set("session_id", sessionId)
		session, err := a.sessionRepository.Get(c, sessionId)
		if err != nil || session ==nil {
				err := &customerror.BusinnesException{
				StatusCode: http.StatusUnauthorized,
				ErrorCode:  customerror.ErrSessionExpired,
				Message:    "Oturumunuz bitmiştir.",
			}
			c.Error(err) // Hatayı Gin'in listesine ekle
			c.Abort()    // İsteği durdur (Handler'a gitmesin)
			return
		}
		c.Set("Role", user.Role(session.Role))
		c.Set("UserId", session.UserId)
		c.Set("PlayerId", session.PlayerId)
		c.Next()
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.GetString("session_id")
		log.Printf("SessionId= %s ", sessionId)

		if sessionId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// Context'e user bilgilerini koy
		//Redisden kullanıcı bilgileri çekilecek
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
				StatusCode: http.StatusUnauthorized,
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
			StatusCode: http.StatusUnauthorized,
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

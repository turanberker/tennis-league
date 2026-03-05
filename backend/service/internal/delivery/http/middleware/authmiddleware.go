package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/turanberker/tennis-league-service/internal/delivery"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.PrivateClaims()
		sessionId := claims["session_id"].(string)
		// Context'e user bilgilerini koy
		c.Set("session_id", sessionId)
		session, err := a.sessionRepository.Get(c, sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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

		roleValue, exists := c.Get("Role")

		roleStrings := make([]string, len(roles))
		for i, r := range roles {
			roleStrings[i] = string(r)
		}

		if !exists {

			c.AbortWithStatusJSON(
				http.StatusForbidden,
				delivery.NewErrorResponse(
					"Kullanıcı şu rollerden birine sahip olmalı: "+strings.Join(roleStrings, ", "),
				),
			)
			return
		}

		userRole, ok := roleValue.(user.Role)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				delivery.NewErrorResponse("Geçersiz rol tipi"),
			)
			return
		}

		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, delivery.NewErrorResponse(
			fmt.Sprintf("Kullanıcı %v rollerinden birine sahip olmalı", strings.Join(roleStrings, ", "))))
	}
}

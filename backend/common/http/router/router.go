package router

import (
	"time"

	"tennis-league/common/lib/error/handler"
	authmiddleware "tennis-league/common/security/authmiddleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(serverConfig *ServerConfig, auth *authmiddleware.AuthMiddleware,
	handlers ...RegisterableHandler) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{serverConfig.AllowedOrigins}, // frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(handler.ErrorHandler())
	r.Use(auth.GetToken()) // 🔥 Global JWT kontrolü

	for _, h := range handlers {
		h.RegisterRoutes(r)
	}

	return r
}

type RegisterableHandler interface {
	RegisterRoutes(r *gin.Engine)
}

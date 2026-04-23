package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	"github.com/turanberker/tennis-league-service/internal/platform"
)

func NewRouter(serverConfig *platform.ServerConfig, auth *middleware.AuthMiddleware,
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
	r.Use(middleware.ErrorHandler())
	r.Use(auth.GetToken()) // 🔥 Global JWT kontrolü

	for _, h := range handlers {
		h.RegisterRoutes(r)
	}

	return r
}

type RegisterableHandler interface {
	RegisterRoutes(r *gin.Engine)
}

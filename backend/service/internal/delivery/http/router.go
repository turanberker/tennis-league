package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handlers ...RegisterableHandler) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	for _, h := range handlers {
		h.RegisterRoutes(r)
	}

	return r
}

type RegisterableHandler interface {
	RegisterRoutes(r *gin.Engine)
}

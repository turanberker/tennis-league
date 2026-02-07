package http

import (
	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/userhandler"
)

func NewRouter(userHandler *userhandler.UserHandler) *gin.Engine {
	r := gin.Default()

	userHandler.RegisterRoutes(r)

	return r
}

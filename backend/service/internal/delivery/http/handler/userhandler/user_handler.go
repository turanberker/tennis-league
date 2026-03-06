package userhandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
)

type UserHandler struct {
	userUc *user.Usecase
}

func NewUserHandler(userUc *user.Usecase) *UserHandler {
	return &UserHandler{userUc: userUc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/user")
	{
		leagues.GET("/list", middleware.RequireRole(user.RoleAdmin), h.getAll)
	}
}

func (h *UserHandler) getAll(c *gin.Context) {
	users, error := h.userUc.GetAll(c.Request.Context())
	if error != nil {
		log.Fatal(error.Error())
		c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
		return
	}
	usersResponse := make([]*UserResponse, 0, len(users))

	for _, l := range users {
		usersResponse = append(usersResponse, toPlayerResponse(l))
	}

	res := delivery.NewSuccessResponse(usersResponse)
	c.JSON(http.StatusOK, res)

}

func toPlayerResponse(l *user.User) *UserResponse {
	if l == nil {
		return nil
	}

	return &UserResponse{
		Id:       l.Id,
		Name:     l.Name,
		Surname:  l.Surname,
		Role:     l.Role,
		Email:    l.Email,
		Approved: l.Approved,
		PlayerId: l.PlayerId,
	}
}

package userhandler

import (
	"net/http"

	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	authmiddleware "tennis-league/common/security/auth"
	"tennis-league/common/security/dto"

	"tennis-league/service/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUc *user.Usecase
}

func NewUserHandler(userUc *user.Usecase) *UserHandler {
	return &UserHandler{userUc: userUc}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {

	userRoute := r.Group("/user")
	{
		userRoute.GET("/list", authmiddleware.RequireRole(dto.RoleAdmin), h.getAll)

		profile := userRoute.Group("/profile", authmiddleware.RequireAuth())
		{
			profile.PATCH("/change-password", h.changeMyPassword)
		}

	}
}

func (h *UserHandler) getAll(c *gin.Context) {
	users, error := h.userUc.GetAll(c.Request.Context())
	if error != nil {
		c.Error(error)
		c.Abort()
		return

	}
	usersResponse := make([]*UserResponse, 0, len(users))

	for _, l := range users {
		usersResponse = append(usersResponse, toPlayerResponse(l))
	}

	res := delivery.NewSuccessResponse(usersResponse)
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) changeMyPassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(customerror.NewValidationError(ve))
			c.Abort()
			return
		} else {
			c.Error(customerror.NewInternalError(err))
			c.Abort()
			return
		}
	}

	userId, _ := authmiddleware.GetUserIdFromContext(c)

	err := h.userUc.ChangePassword(c.Request.Context(), userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	res := delivery.NewSuccessResponse("password changed successfully")
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

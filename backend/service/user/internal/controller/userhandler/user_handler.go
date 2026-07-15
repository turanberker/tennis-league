package userhandler

import (
	"tennis-league/common/http/router"
	"tennis-league/user-service/internal/service/user"

	authmiddleware "tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"
)

type UserHandler struct {
	userUc *user.Usecase
}

func NewUserHandler(userUc *user.Usecase) *UserHandler {
	return &UserHandler{userUc: userUc}
}

func (h *UserHandler) RegisterRoutes(r *router.CustomRouterGroup) {

	userRoute := r.Group("/user")
	{
		userRoute.GET("/list", authmiddleware.RequireRole(dto.RoleAdmin), h.getAll)

		profile := userRoute.Group("/profile", authmiddleware.RequireAuth())
		{
			profile.PATCH("/change-password", h.changeMyPassword)
		}

	}
}

func (h *UserHandler) getAll(c *router.CustomContext) {
	users, err := h.userUc.GetAll(c.Request.Context())
	if err != nil {
		c.ErrorComplete(err)
		return
	}
	usersResponse := make([]*UserResponse, 0, len(users))

	for _, l := range users {
		usersResponse = append(usersResponse, toPlayerResponse(l))
	}
	c.OkComplete(usersResponse)
}

func (h *UserHandler) changeMyPassword(c *router.CustomContext) {
	var req struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword"`
	}

	if !c.BindJSONOrAbort(&req) {
		return
	}

	userId, _ := c.CurrentUserId()

	err := h.userUc.ChangePassword(c.Request.Context(), userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.ErrorComplete(err)
		return
	}
	c.OkComplete("password changed successfully")
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

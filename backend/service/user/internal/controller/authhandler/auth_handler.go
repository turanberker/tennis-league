package authhandler

import (
	"net/http"
	"tennis-league/common/http/router"
	service "tennis-league/user-service/internal"
	"tennis-league/user-service/internal/service/auth"
	"tennis-league/user-service/internal/service/token"

	customerror "tennis-league/common/lib/error"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc           *auth.Usecase
	tokenService *token.TokenService
}

func NewAuthHandler(uc *auth.Usecase, tokenService *token.TokenService) *AuthHandler {
	return &AuthHandler{uc: uc, tokenService: tokenService}
}

func (h *AuthHandler) RegisterRoutes(r *router.CustomRouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", h.login)
		authGroup.POST("/refresh", h.refresh)
		authGroup.POST("/register", h.register)
		authGroup.POST("/logout", h.logout)
	}
}

func (h *AuthHandler) logout(c *router.CustomContext) {
	sId, _ := c.Get("session_id")

	// 1. Redis'ten sil
	if sessionId, ok := sId.(string); ok {
		h.uc.DeleteSessionFromRedis(c.Request.Context(), sessionId)
	}

	// 2. Tarayıcıdaki cookie'leri temizle (Sürelerini -1 yaparak)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/authmiddleware/refresh", "", false, true)
	c.OkComplete("Çıkış Yaptınız")
}

func (h *AuthHandler) login(c *router.CustomContext) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if !c.BindJSONOrAbort(&req) {
		return
	}

	usr, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.ErrorComplete(customerror.NewBusinessError(http.StatusUnauthorized,
			service.INVALID_CREDENTIAL, "invalid email or password"))
		return
	}

	_, err = h.tokenService.GenerateAccessTokenAndSetCookie(c, usr.SessionId)

	if err != nil {
		c.ErrorComplete(err)
		return
	}

	_, err = h.tokenService.GenerateRefreshTokenAndSetCookie(c, usr.SessionId)

	if err != nil {
		c.ErrorComplete(err)
		return
	}

	c.OkComplete(LoginResponse{
		CurrentUser: CurrentUserDTO{
			UserID:   usr.ID,
			Name:     usr.Name,
			Surname:  usr.Surname,
			Role:     string(usr.Role),
			PlayerId: usr.PlayerId,
		},
	})
}

func (h *AuthHandler) refresh(c *router.CustomContext) {
	// Cookie'den refresh token'ı oku
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Token'ı valide et ve yeni bir Access Token üret
	newAccessToken, err := h.tokenService.ValidateAndRefreshAndSetAccessCookie(c, refreshToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newAccessToken})
}

func (h *AuthHandler) register(c *router.CustomContext) {
	var req RegisterRequest

	if !c.BindJSONOrAbort(&req) {
		return
	}

	usr, err := h.uc.RegisterUser(
		c.Request.Context(),
		&auth.RegisterUserInput{
			Email:    req.Email,
			Name:     req.Name,
			Surname:  req.Surname,
			Password: req.Password,
		},
	)

	if err != nil {
		c.ErrorComplete(err)
		return
	}

	_, err = h.tokenService.GenerateAccessTokenAndSetCookie(c, usr.SessionId)
	if err != nil {
		c.ErrorComplete(err)
		return
	}

	_, err = h.tokenService.GenerateRefreshTokenAndSetCookie(c, usr.SessionId)
	if err != nil {
		c.ErrorComplete(err)
		return
	}

	// JWT oluştur

	c.OkComplete(LoginResponse{
		CurrentUser: CurrentUserDTO{
			UserID:  usr.ID,
			Name:    usr.Name,
			Surname: usr.Surname,
			Role:    string(usr.Role),
		},
	})
}

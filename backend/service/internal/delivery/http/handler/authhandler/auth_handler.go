package authhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"

	"github.com/turanberker/tennis-league-service/internal/domain/auth"
	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
)

type AuthHandler struct {
	uc           *auth.Usecase
	tokenService *middleware.TokenService
}

func NewAuthHandler(uc *auth.Usecase, tokenService *middleware.TokenService) *AuthHandler {
	return &AuthHandler{uc: uc, tokenService: tokenService}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", h.login)
		authGroup.POST("/refresh", h.refresh)
		authGroup.POST("/register", h.register)
		authGroup.POST("/logout", h.logout)
	}
}

func (h *AuthHandler) logout(c *gin.Context) {
	sId, _ := c.Get("session_id")

	// 1. Redis'ten sil
	if sessionId, ok := sId.(string); ok {
		h.uc.DeleteSessionFromRedis(c.Request.Context(), sessionId)
	}

	// 2. Tarayıcıdaki cookie'leri temizle (Sürelerini -1 yaparak)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", false, true)

	c.JSON(http.StatusOK, delivery.NewSuccessResponse("Çıkış Yaptınız"))
}

func (h *AuthHandler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
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

	usr, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.Error(customerror.NewBussinnessError(http.StatusUnauthorized,
			customerror.INVALID_CREDENTIAL, "invalid email or password"))
		c.Abort()
		return
	}

	_, err = h.tokenService.GenerateAccessTokenAndSetCookie(c, usr.SessionId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	_, err = h.tokenService.GenerateRefreshTokenAndSetCookie(c, usr.SessionId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	response := delivery.NewSuccessResponse(LoginResponse{
		CurrentUser: CurrentUserDTO{
			UserID:   usr.ID,
			Name:     usr.Name,
			Surname:  usr.Surname,
			Role:     string(usr.Role),
			PlayerId: usr.PlayerId,
		},
	})

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) refresh(c *gin.Context) {
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

func (h *AuthHandler) register(c *gin.Context) {
	var req RegisterRequest

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

	usr, err := h.uc.RegisterUser(
		c.Request.Context(),
		&user.RegisterUserInput{
			Email:    req.Email,
			Name:     req.Name,
			Surname:  req.Surname,
			Password: req.Password,
		},
	)

	if err != nil {
		var be *customerror.BusinnesException
		if errors.As(err, &be) {
			c.Error(be)
			c.Abort()
			return
		} else {
			c.Error(customerror.NewInternalError(err))
			c.Abort()
			return
		}
	}

	_, err = h.tokenService.GenerateAccessTokenAndSetCookie(c, usr.SessionId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	_, err = h.tokenService.GenerateRefreshTokenAndSetCookie(c, usr.SessionId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// JWT oluştur
	response := delivery.NewSuccessResponse(LoginResponse{
		CurrentUser: CurrentUserDTO{
			UserID:  usr.ID,
			Name:    usr.Name,
			Surname: usr.Surname,
			Role:    string(usr.Role),
		},
	})

	c.JSON(http.StatusOK, response)
}

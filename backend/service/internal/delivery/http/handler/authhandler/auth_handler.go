package authhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	r.POST("auth/login", h.login)
	r.POST("auth/register", h.register)
}

func (h *AuthHandler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, delivery.NewErrorResponse(err.Error()))
		return
	}

	usr, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.Error(customerror.NewBussinnessError(http.StatusUnauthorized,
			 customerror.INVALID_CREDENTIAL, "invalid email or password"))
		c.Abort()
		return
	}

	tokenString, _ := h.tokenService.Generate(usr.SessionId)
	// JWT oluştur

	response := delivery.NewSuccessResponse(LoginResponse{
		Token: tokenString,
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

func (h *AuthHandler) register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, _ := h.tokenService.Generate(usr.SessionId)
	// JWT oluştur
	response := delivery.NewSuccessResponse(LoginResponse{
		Token: tokenString,
		CurrentUser: CurrentUserDTO{
			UserID:  usr.ID,
			Name:    usr.Name,
			Surname: usr.Surname,
			Role:    string(usr.Role),
		},
	})

	c.JSON(http.StatusOK, response)
}

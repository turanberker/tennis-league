package userhandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
)

type UserHandler struct {
	uc        *user.Usecase
	tokenAuth *jwtauth.JWTAuth
}

func NewUserHandler(uc *user.Usecase, tokenAuth *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{uc: uc, tokenAuth: tokenAuth}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("auth/login", h.login)
	r.POST("auth/register", h.register)
}

func (h *UserHandler) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	tokenString, _ := h.createToken(usr)
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

func (h *UserHandler) createToken(usr *user.User) (string, error) {
	_, tokenString, err := h.tokenAuth.Encode(map[string]interface{}{
		"session_id": uuid.New(),
		"user_id":    usr.ID,
		"role":       usr.Role,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	})
	return tokenString, err
}

func (h *UserHandler) register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	tokenString, _ := h.createToken(usr)
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

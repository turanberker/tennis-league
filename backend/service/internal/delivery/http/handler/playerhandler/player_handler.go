package playerhandler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type PlayerHandler struct {
	tm *database.TransactionManager
	uc *player.Usecase
}

func NewPlayerHandler(uc *player.Usecase, tm *database.TransactionManager) *PlayerHandler {
	return &PlayerHandler{uc: uc, tm: tm}
}

func (h *PlayerHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/player")
	{
		group.GET("/list", h.getAll)
		group.POST("", middleware.RequireRole(user.RoleAdmin), h.save)
		group.PUT("/:id", middleware.RequireRole(user.RoleAdmin), h.assignToUser)
		group.GET("/:uuid", h.getByUuid)
		group.GET("/unassigned-players", h.unassignedPlayers)
	}

}

func (h *PlayerHandler) getByUuid(c *gin.Context) {
	uuidParam := c.Param("uuid")
	if uuidParam == "" {
		c.JSON(http.StatusBadRequest, delivery.NewErrorResponse("Uuid is required"))
		return
	}
	player, err := h.uc.GetByUuid(c.Request.Context(), uuidParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, delivery.NewErrorResponse("player is not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, delivery.NewErrorResponse("internal error"))

		return
	}
	response := toPlayerResponse(player)
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(response))

}

func (h *PlayerHandler) save(c *gin.Context) {

	var req struct {
		Name    string     `json:"name" binding:"min=3,max=75,required"`
		Surname string     `json:"surname" binding:"min=3,max=75,required"`
		Sex     player.Sex `json:"sex" binding:"required,oneof=M F"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persistPlayer := &player.PersistPlayer{
		Name:    req.Name,
		Surname: req.Surname,
		Sex:     req.Sex,
	}

	playerId, err := h.uc.Save(c.Request.Context(), persistPlayer)

	if err != nil {
		res := delivery.NewErrorResponse(err.Error())
		c.JSON(http.StatusOK, res)
	} else {
		res := delivery.NewSuccessResponse(playerId)
		c.JSON(http.StatusOK, res)
	}
}

func (h *PlayerHandler) assignToUser(c *gin.Context) {
	playerId := c.Param("id")

	var req struct {
		UserId string `form:"userId" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}
	var err error
	ctx := c.Request.Context()

	ctx, err = h.tm.StartTransaction(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	defer database.DeferRollback(ctx, &err)

	err = h.uc.AssignToUser(ctx, playerId, req.UserId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	err = database.Commit(ctx)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse("İşlem başarılı"))
}

func (h *PlayerHandler) getAll(c *gin.Context) {
	var req struct {
		Name *string     `form:"name" binding:"omitempty"`
		Sex  *player.Sex `form:"sex" binding:"omitempty,oneof=M F"`
	}

	// Gin otomatik olarak URL'deki ?name=...&sex=... kısımlarını struct'a doldurur
	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	players, _ := h.uc.List(c.Request.Context(), player.ListQueryParameters{Name: req.Name,
		Sex: req.Sex})

	playersResponse := make([]*PlayerResponse, 0, len(players))

	for _, l := range players {
		playersResponse = append(playersResponse, toPlayerResponse(l))
	}

	res := delivery.NewSuccessResponse(playersResponse)
	c.JSON(http.StatusOK, res)
}

func (h *PlayerHandler) unassignedPlayers(c *gin.Context) {

	var req struct {
		Sex player.Sex `form:"sex" binding:"oneof=M F"`
	}
	// Gin otomatik olarak URL'deki ?name=...&sex=... kısımlarını struct'a doldurur
	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	isFalse := false
	players, err := h.uc.List(c.Request.Context(),
		player.ListQueryParameters{Sex: &req.Sex,
			HasUser: &isFalse})

	if err != nil {
		c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
		return
	}

	leagueResponse := make([]*PlayerResponse, 0, len(players))

	for _, l := range players {
		leagueResponse = append(leagueResponse, toPlayerResponse(l))
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(leagueResponse))

}

func toPlayerResponse(l *player.Player) *PlayerResponse {
	if l == nil {
		return nil
	}

	return &PlayerResponse{
		ID:      l.ID,
		Name:    l.Name,
		Surname: l.Surname,
		Sex:     l.Sex,
		UserId:  l.UserId,
	}
}

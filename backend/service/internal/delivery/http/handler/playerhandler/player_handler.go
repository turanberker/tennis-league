package playerhandler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
)

type PlayerHandler struct {
	uc *player.Usecase
}

func NewPlayerHandler(uc *player.Usecase) *PlayerHandler {
	return &PlayerHandler{uc: uc}
}

func (h *PlayerHandler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/player")
	{
		leagues.GET("/list", h.getAll)
		leagues.POST("", h.save)
		leagues.GET("/:uuid", h.getByUuid)
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

func (h *PlayerHandler) getAll(c *gin.Context) {
	name := c.Query("name") // query param
	players, _ := h.uc.List(c.Request.Context(), name)

	playersResponse := make([]*PlayerResponse, 0, len(players))

	for _, l := range players {
		playersResponse = append(playersResponse, toPlayerResponse(l))
	}

	res := delivery.NewSuccessResponse(playersResponse)
	c.JSON(http.StatusOK, res)
}

func toPlayerResponse(l *player.Player) *PlayerResponse {
	if l == nil {
		return nil
	}

	return &PlayerResponse{
		ID:      l.ID,
		Name:    l.Name,
		Surname: l.Surname,
		UUID:    l.Uuid,
		Sex:     l.Sex,
		UserId:  l.UserId,
	}
}

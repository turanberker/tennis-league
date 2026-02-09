package leaguehandler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
)

type Handler struct {
	uc *league.Usecase
}

func NewHandler(uc *league.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/leagues")
	{
		leagues.GET("/list", h.getAll)
		leagues.POST("", h.save)
		leagues.GET("/:id", h.getById)
	}

}

func (h *Handler) getById(c *gin.Context) {
	ctx := c.Request.Context()

	// path param
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.NewErrorResponse("invalid id"))

		return
	}

	league, err := h.uc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, delivery.NewErrorResponse("league not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, delivery.NewErrorResponse("internal error"))

		return
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(league))
}

func (h *Handler) save(c *gin.Context) {

	var req struct {
		Name string `json:"name" binding:"min=3,max=75,required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persistLeague := &league.PersistLeague{
		Name: req.Name,
	}

	leagueId, err := h.uc.Save(c.Request.Context(), persistLeague)

	if err != nil {
		res := delivery.NewErrorResponse(err.Error())
		c.JSON(http.StatusOK, res)
	} else {
		res := delivery.NewSuccessResponse(leagueId)
		c.JSON(http.StatusOK, res)
	}

}

func (h *Handler) getAll(c *gin.Context) {

	name := c.Query("name") // query param
	leagues, _ := h.uc.GetAll(c.Request.Context(), name)

	res := delivery.NewSuccessResponse(leagues)
	c.JSON(http.StatusOK, res)

}

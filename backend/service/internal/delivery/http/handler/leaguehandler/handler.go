package leaguehandler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/dto"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
	"github.com/turanberker/tennis-league-service/internal/domain/team"
)

type Handler struct {
	uc     *league.Usecase
	teamUc *team.UseCase
}

func NewHandler(uc *league.Usecase, teamUc *team.UseCase) *Handler {
	return &Handler{uc: uc,
		teamUc: teamUc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/leagues")
	{
		leagues.GET("/list", h.getAll)
		leagues.POST("", h.save)
		leagues.GET("/:id", h.getById)
		leagues.POST("/:id/create-fixture", h.createFixture)
		leagues.GET("/:id/teams", h.getTeams)
		leagues.POST("/:id/teams", h.newTeam)
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
	response := toLeagueResponse(league)
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(response))
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

	leagueResponse := make([]*LeagueResponse, 0, len(leagues))

	for _, l := range leagues {
		leagueResponse = append(leagueResponse, toLeagueResponse(l))
	}

	res := delivery.NewSuccessResponse(leagueResponse)
	c.JSON(http.StatusOK, res)

}

func (h *Handler) getTeams(c *gin.Context) {

	idParam := c.Param("id") // query param

	teams, err := h.teamUc.GetByLeagueId(c.Request.Context(), idParam)

	if err != nil {
		res := delivery.NewErrorResponse("Takımlar Çekilemedi")
		c.JSON(http.StatusOK, res)
		return
	}

	teamResponse := make([]*dto.TeamResponse, 0, len(teams))

	for _, l := range teams {
		teamResponse = append(teamResponse, toTeamResponse(l))
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(teamResponse))
}

func (h *Handler) newTeam(c *gin.Context) {

	leagueId := c.Param("id") // query param

	var req struct {
		Name      string   `json:"name" binding:"min=3,max=75,required"`
		PlayerIDs []string `json:"playerIds" binding:"required,len=2,dive,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.teamUc.Save(c.Request.Context(), &team.CreateTeamRequest{
		LeagueID:  leagueId,
		Name:      req.Name,
		PlayerIDs: req.PlayerIDs,
	})
	if err != nil {

		res := delivery.NewErrorResponse(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	res := delivery.NewSuccessResponse(id)
	c.JSON(http.StatusOK, res)

}

func (h *Handler) createFixture(c *gin.Context) {

	leagueId := c.Param("id") // query param
	h.uc.SetFitxtureCreatedDate(c.Request.Context(), leagueId)

	res := delivery.NewSuccessResponse("Fikstür oluşturuldu")
	c.JSON(http.StatusOK, res)
}

func toLeagueResponse(l *league.League) *LeagueResponse {
	if l == nil {
		return nil
	}

	return &LeagueResponse{
		ID:                 l.ID,
		Name:               l.Name,
		FixtureCreatedDate: l.FixtureCreatedDate,
	}
}

func toTeamResponse(l *team.Team) *dto.TeamResponse {
	if l == nil {
		return nil
	}

	return &dto.TeamResponse{
		ID:   l.ID,
		Name: l.Name,
	}
}

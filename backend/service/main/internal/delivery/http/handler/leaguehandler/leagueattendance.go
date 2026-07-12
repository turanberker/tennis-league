package leaguehandler

import (
	"net/http"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/team"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type leagueAttendanceHandler struct {
	leagueHandlerMiddleware *leagueHandlerMiddleware
	teamUc                  *team.UseCase
	uc                      *league.Usecase
}

func (h *leagueAttendanceHandler) registerSubRoutes(group *gin.RouterGroup) {
	group.GET("/teams", h.getTeams)
	group.POST("/teams",
		authmiddleware.RequireRole(dto.RoleAdmin, dto.RoleCoordinator),
		h.leagueHandlerMiddleware.checkIfCoordinator,
		h.newTeam)
	group.GET("/players", h.players)
	group.POST("/players", h.addPlayer)
}

func (h *leagueAttendanceHandler) getTeams(c *gin.Context) {

	idParam := c.Param("id") // query param

	teams, err := h.teamUc.GetByLeagueId(c.Request.Context(), idParam)

	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	type TeamResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Power int32  `json:"power"`
	}

	teamResponse := make([]*TeamResponse, 0, len(teams))

	for _, l := range teams {

		res := &TeamResponse{
			ID:    l.ID,
			Name:  l.Name,
			Power: l.Power,
		}
		teamResponse = append(teamResponse, res)
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(teamResponse))
}

func (h *leagueAttendanceHandler) newTeam(c *gin.Context) {

	leagueId := c.Param("id") // query param

	var req struct {
		Name      string   `json:"name" binding:"min=3,max=75,required"`
		PlayerIDs []string `json:"playerIds" binding:"required,len=2,dive,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			_ = c.Error(customerror.NewValidationError(ve))
			c.Abort()
			return
		} else {
			_ = c.Error(customerror.NewInternalError(err))
			c.Abort()
			return
		}
	}

	response, err := h.uc.CreateTeam(c.Request.Context(), &league.CreateTeamRequestDto{
		LeagueId:  leagueId,
		Name:      req.Name,
		PlayerIDs: req.PlayerIDs,
	})

	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	var resModel struct {
		TeamId               string `json:"teamId"`
		TotalAttendanceCount int32  `json:"totalAttendanceCount"`
	}
	resModel.TeamId = response.TeamId
	resModel.TotalAttendanceCount = response.TotalAttendance
	res := delivery.NewSuccessResponse(resModel)
	c.JSON(http.StatusOK, res)

}

func (h *leagueAttendanceHandler) players(c *gin.Context) {
	idParam := c.Param("id")
	players, err := h.uc.GetPlayersByLeagueId(c.Request.Context(), idParam)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	type PlayerResponse struct {
		ID        string `json:"id"`
		FirstName string `json:"firstname"`
		SurName   string `json:"surname"`
		Power     int    `json:"power"`
	}
	response := make([]PlayerResponse, 0, len(players))

	for _, l := range players {

		pr := PlayerResponse{
			ID:        l.ID,
			FirstName: l.Firstname,
			SurName:   l.Surname,
			Power:     l.Power,
		}
		response = append(response, pr)
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(response))

}

func (h *leagueAttendanceHandler) addPlayer(c *gin.Context) {
	leagueId := c.Param("id")

	var req struct {
		PlayerId string `form:"playerId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	totalAttendance, err := h.uc.AddPlayerToLeague(c.Request.Context(), leagueId, req.PlayerId)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	type addPlayerResponse struct {
		PlayerId             string `json:"playerId"`
		TotalAttendanceCount *int32 `json:"totalAttendanceCount"`
	}
	response := addPlayerResponse{PlayerId: req.PlayerId, TotalAttendanceCount: totalAttendance}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(response))
}

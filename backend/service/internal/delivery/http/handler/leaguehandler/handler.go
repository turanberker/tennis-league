package leaguehandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
	"github.com/turanberker/tennis-league-service/internal/domain/team"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type Handler struct {
	tm           *database.TransactionManager
	uc           *league.Usecase
	teamUc       *team.UseCase
	scoreBaordUc *scoreboard.UseCase
}

func NewHandler(uc *league.Usecase, teamUc *team.UseCase, scoreBaordUc *scoreboard.UseCase) *Handler {
	return &Handler{uc: uc,
		teamUc:       teamUc,
		scoreBaordUc: scoreBaordUc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/leagues")
	{
		leagues.GET("/list", h.getAll)
		leagues.POST("", middleware.RequireRole(user.RoleAdmin), h.save)
		leagues.GET("/:id", h.getById)
		leagues.POST("/:id/create-fixture",
			middleware.RequireRole(user.RoleAdmin, user.RoleCoordinator),
			h.checkIfCoordinator,
			h.createFixture)
		leagues.GET("/:id/teams", h.getTeams)
		leagues.POST("/:id/teams",
			middleware.RequireRole(user.RoleAdmin, user.RoleCoordinator),
			h.checkIfCoordinator,
			h.newTeam)
		leagues.GET("/:id/fixture", h.getFixture)
		leagues.GET("/:id/standings", h.getScoreBoard)
		leagues.POST("/:id/coordinator",
			middleware.RequireRole(user.RoleAdmin, user.RoleCoordinator),
			h.checkIfCoordinator, h.newCoordinator)
	}

}

func (h *Handler) checkIfCoordinator(c *gin.Context) {
	roleValue, _ := c.Get("Role")
	leagueId := c.Param("id")
	userIdValue, _ := c.Get("UserId")
	userId, _ := userIdValue.(string)

	if role, ok := roleValue.(user.Role); ok {

		// 3. Karşılaştırma yap
		if role == user.RoleCoordinator {
			coordinator, err := h.uc.IsUserCoordinator(c.Request.Context(), leagueId, userId)
			if err != nil {
				c.Error(customerror.NewInternalError(err))
				c.Abort()
			}
			if coordinator {
				c.Next()
			} else {
				err := &customerror.BusinnesException{
					StatusCode: http.StatusForbidden,
					ErrorCode:  customerror.INSUFFICIENT_PERMISSIONS,
					Message:    "Bu ligde koordinatör değilsiniz",
				}
				c.Error(err)
				c.Abort()
			}
		}

		if role == user.RoleAdmin {
			c.Next()
		}
	} else {
		err := &customerror.BusinnesException{
			StatusCode: http.StatusForbidden,
			ErrorCode:  customerror.INSUFFICIENT_PERMISSIONS,
			Message:    "Bu ligde yetkiniz yok",
		}
		c.Error(err)
		c.Abort()
	}
}

func (h *Handler) getById(c *gin.Context) {
	ctx := c.Request.Context()

	// path param
	leagueId := c.Param("id")

	leagueData, err := h.uc.GetById(ctx, leagueId)
	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	var res struct {
		Id              string                     `json:"id"`
		Name            string                     `json:"name"`
		Format          league.LEAGUE_FORMAT       `json:"format"`
		Category        league.LEAGUE_CATEGORY     `json:"category"`
		ProcessType     league.LEAGUE_PROCESS_TYPE `json:"processType"`
		Status          league.LEAGUE_STATUS       `json:"status"`
		TotalAttentance int32                      `json:"totalAttentance"`
		StartedDate     *time.Time                 `json:"startedDate"`
		EndDate         *time.Time                 `json:"endDate"`
	}

	// Alanları doldur
	res.Id = leagueData.ID
	res.Name = leagueData.Name
	res.Format = leagueData.Format
	res.Category = leagueData.Category
	res.ProcessType = leagueData.ProcessType
	res.Status = leagueData.Status
	res.TotalAttentance = leagueData.TotalAttendance
	res.StartedDate = leagueData.StartDate
	res.EndDate = leagueData.EndDate

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(res))
}

func (h *Handler) save(c *gin.Context) {

	var req struct {
		Name        string                     `json:"name" binding:"min=3,max=75,required"`
		Format      league.LEAGUE_FORMAT       `json:"format" binding:"required"`
		Categoty    league.LEAGUE_CATEGORY     `json:"category" binding:"required"`
		ProcessType league.LEAGUE_PROCESS_TYPE `json:"processType" binding:"required"`
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

	persistLeague := &league.PersistLeague{
		Name:        req.Name,
		Format:      req.Format,
		Categoty:    req.Categoty,
		ProcessType: req.ProcessType,
	}

	leagueId, err := h.uc.Save(c.Request.Context(), persistLeague)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	} else {
		res := delivery.NewSuccessResponse(leagueId)
		c.JSON(http.StatusOK, res)
	}

}

func (h *Handler) getAll(c *gin.Context) {

	var req struct {
		Status *league.LEAGUE_STATUS `form:"status" binding:"omitempty"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
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

	leagues, err := h.uc.GetAll(c.Request.Context(), req.Status)
	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	type response struct {
		ID                 string                     `json:"id"`
		Name               string                     `json:"name"`
		Format             league.LEAGUE_FORMAT       `json:"format"`
		Category           league.LEAGUE_CATEGORY     `json:"category"`
		ProcessType        league.LEAGUE_PROCESS_TYPE `json:"processType"`
		Status             league.LEAGUE_STATUS       `json:"status"`
		TotalAttentance    int32                      `json:"totalAttentance"`
		CoordinatorUserIds []string                   `json:"coordinatorUserIds"`
	}
	leagueResponse := make([]*response, 0, len(leagues))

	for _, l := range leagues {
		leagueResponse = append(leagueResponse, &response{
			ID:                 l.ID,
			Name:               l.Name,
			Format:             l.Format,
			Category:           l.Category,
			ProcessType:        l.Type,
			Status:             l.Status,
			TotalAttentance:    l.TotalAttentance,
			CoordinatorUserIds: l.CoordinatorUserId,
		})
	}

	res := delivery.NewSuccessResponse(leagueResponse)
	c.JSON(http.StatusOK, res)

}

func (h *Handler) getTeams(c *gin.Context) {

	idParam := c.Param("id") // query param

	teams, err := h.teamUc.GetByLeagueId(c.Request.Context(), idParam)

	if err != nil {
		c.Error(err)
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

		team := &TeamResponse{
			ID:    l.ID,
			Name:  l.Name,
			Power: l.Power,
		}
		teamResponse = append(teamResponse, team)
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

	response, err := h.uc.CreateTeam(c.Request.Context(), &league.CreateTeamRequestDto{
		LeagueId:  leagueId,
		Name:      req.Name,
		PlayerIDs: req.PlayerIDs,
	})

	if err != nil {
		c.Error(err)
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

func (h *Handler) createFixture(c *gin.Context) {

	leagueId := c.Param("id") // query param
	err := h.uc.CreateFixture(c.Request.Context(), leagueId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	res := delivery.NewSuccessResponse("Fikstür oluşturuldu")
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getFixture(c *gin.Context) {
	leagueId := c.Param("id") // query param

	fixture, err := h.uc.GetFixture(c.Request.Context(), leagueId)

	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}
	fixtureResponse := make([]*LeagueFixtureMatchResponse, 0, len(fixture))

	for _, l := range fixture {
		fixtureResponse = append(fixtureResponse, toFixtureResponse(l))
	}
	res := delivery.NewSuccessResponse(fixtureResponse)
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getScoreBoard(c *gin.Context) {
	leagueId := c.Param("id")

	board, err := h.scoreBaordUc.GetScoreBoard(c.Request.Context(), leagueId)
	if err != nil {
		res := delivery.NewErrorResponse(err.Error())
		c.JSON(http.StatusOK, res)
		return
	}

	var result []*ScoreBoardResponse
	for o, b := range board {
		team := &ScoreBoardResponse{
			TeamRef:   TeamRef{Id: b.Team.Id, Name: b.Team.Name},
			Order:     o + 1,
			Played:    b.Played,
			Won:       b.Won,
			Lost:      b.Lost,
			WonSets:   b.WonSets,
			LostSets:  b.LostSets,
			WonGames:  b.WonGames,
			LostGames: b.LostGames,
			Score:     b.Score,
		}
		result = append(result, team)
	}
	res := delivery.NewSuccessResponse(result)
	c.JSON(http.StatusOK, res)

}

func (h *Handler) newCoordinator(c *gin.Context) {
	leagueId := c.Param("id")
	var req struct {
		UserId string `form:"userId" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
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

	added, err := h.uc.AddNewCoordinator(c.Request.Context(), leagueId, req.UserId)

	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	res := delivery.NewSuccessResponse(added)
	c.JSON(http.StatusOK, res)

}

func toFixtureResponse(l *match.LeagueFixtureMatch) *LeagueFixtureMatchResponse {
	if l == nil {
		return nil
	}
	return &LeagueFixtureMatchResponse{
		Id:        l.Id,
		Team1:     TeamRefResponse{TeamRef: TeamRef{Id: l.Team1.Id, Name: l.Team1.Name}, Score: l.Team1.Score, Winner: l.Team1.Winner},
		Team2:     TeamRefResponse{TeamRef: TeamRef{Id: l.Team2.Id, Name: l.Team2.Name}, Score: l.Team2.Score, Winner: l.Team2.Winner},
		Status:    l.Status,
		MatchDate: l.MatchDate,
	}
}

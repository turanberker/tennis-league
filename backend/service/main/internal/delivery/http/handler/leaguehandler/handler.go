package leaguehandler

import (
	"log"
	"net/http"
	"time"

	"tennis-league/common/lib/database"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	authmiddleware "tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"

	"tennis-league/service/internal/delivery/http/handler/matchhandler"

	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/match"
	"tennis-league/service/internal/domain/scoreboard"
	"tennis-league/service/internal/domain/team"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	tm                      *database.TransactionManager
	uc                      *league.Usecase
	teamUc                  *team.UseCase
	scoreBoardUc            *scoreboard.UseCase
	matchUc                 *match.UseCase
	leagueHandlerMiddleware *leagueHandlerMiddleware
	leagueAttendanceHandler *leagueAttendanceHandler
}

func NewHandler(uc *league.Usecase, teamUc *team.UseCase, scoreBoardUc *scoreboard.UseCase, matchUc *match.UseCase) *Handler {

	leagueHandlerMiddleware :=
		&leagueHandlerMiddleware{uc: uc, matchUc: matchUc}

	leagueAttendanceHandler := &leagueAttendanceHandler{leagueHandlerMiddleware: leagueHandlerMiddleware, teamUc: teamUc, uc: uc}
	return &Handler{uc: uc,
		teamUc:                  teamUc,
		scoreBoardUc:            scoreBoardUc,
		matchUc:                 matchUc,
		leagueHandlerMiddleware: leagueHandlerMiddleware,
		leagueAttendanceHandler: leagueAttendanceHandler,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	leagues := r.Group("/leagues")
	{
		leagues.GET("/list", h.getAll)
		leagues.POST("", authmiddleware.RequireRole(dto.RoleAdmin), h.save)
		leagues.GET("/:id", h.getById)
		leagues.POST("/:id/start",
			authmiddleware.RequireRole(dto.RoleAdmin, dto.RoleCoordinator),
			h.leagueHandlerMiddleware.checkIfCoordinator,
			h.startLeague)
		attendance := leagues.Group("/:id/attendance")
		h.leagueAttendanceHandler.registerSubRoutes(attendance)

		leagues.GET("/:id/fixture", h.getFixture)
		leagues.GET("/:id/standings", h.getScoreBoard)
		leagues.POST("/:id/coordinator",
			authmiddleware.RequireRole(dto.RoleAdmin, dto.RoleCoordinator),
			h.leagueHandlerMiddleware.checkIfCoordinator, h.newCoordinator)
		leagues.PUT("/:id/match/:matchId/update-score", authmiddleware.RequireAuth(),
			h.leagueHandlerMiddleware.checkIfMatchIsLeague,
			h.leagueHandlerMiddleware.checkIfUserIsCoordinatAdminOrPlayer,
			h.updateScore)
		leagues.PUT("/:id/match/:matchId/update-date",
			authmiddleware.RequireRole(dto.RoleAdmin, dto.RoleCoordinator),
			h.leagueHandlerMiddleware.checkIfCoordinator,
			h.updateMatchDate)
		leagues.PUT("/:id/match/:matchId/approve",
			authmiddleware.RequireRole(dto.RoleAdmin, dto.RoleCoordinator),
			h.leagueHandlerMiddleware.checkIfCoordinator,
			h.approveScore)
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

func (h *Handler) startLeague(c *gin.Context) {

	leagueId := c.Param("id") // query param
	err := h.uc.Start(c.Request.Context(), leagueId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	res := delivery.NewSuccessResponse("Lig Başladı")
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getFixture(c *gin.Context) {
	leagueId := c.Param("id") // query param

	var req struct {
		TeamId *string `form:"teamId" binding:"omitempty"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	var filterParam league.FixtureFilterParam
	if req.TeamId != nil {
		filterParam.TeamId = req.TeamId
	}
	fixture, err := h.uc.GetFixture(c.Request.Context(), leagueId, &filterParam)

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

	board, err := h.scoreBoardUc.GetScoreBoard(c.Request.Context(), leagueId)
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

func (h *Handler) updateMatchDate(c *gin.Context) {
	matchId := c.Param("matchId")

	var req struct {
		MatchDate *time.Time `form:"match-date" time_format:"2006-01-02T15:04:05Z07:00"`
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

	h.matchUc.UpdateMatchDate(c.Request.Context(), matchId, match.MatchSource_LEAGUE, req.MatchDate)
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(true))

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

func (h *Handler) approveScore(c *gin.Context) {
	matchId := c.Param("matchId")
	leagueId := c.Param("id")
	err := h.uc.ApproveMatchScore(c.Request.Context(), leagueId, matchId)
	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(nil))
}

func (h *Handler) updateScore(c *gin.Context) {
	matchId := c.Param("matchId")

	macScore := matchhandler.UpdateScoreRequest{}

	if err := c.ShouldBindJSON(&macScore); err != nil {
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

	log.Printf("match id: %s", matchId)
	log.Printf("score :%+v", macScore)

	set1 := match.SaveScore{Team1Score: macScore.Set1.Team1Score, Team2Score: macScore.Set1.Team2Score}
	set2 := match.SaveScore{Team1Score: macScore.Set2.Team1Score, Team2Score: macScore.Set2.Team2Score}

	saveMatchScore := &match.SaveMatchScore{MatchId: matchId, MatchDate: macScore.MatchDate, Set1: set1, Set2: set2}

	if macScore.SuperTie != nil {

		saveMatchScore.SuperTie = &match.SaveScore{}
		saveMatchScore.SuperTie.Team1Score = macScore.SuperTie.Team1Score
		saveMatchScore.SuperTie.Team2Score = macScore.SuperTie.Team2Score

	}

	response, err := h.matchUc.SaveMatchScore(c.Request.Context(), saveMatchScore)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(&matchhandler.MatchScoreResponse{Team1Score: response.Team1Score, Team2Score: response.Team2Score}))

}

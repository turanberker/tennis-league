package dashboard

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
)

type DashboardHandler struct {
	playerUc *player.Usecase
}

func NewDashboardHandler(playerUc *player.Usecase) *DashboardHandler {
	return &DashboardHandler{playerUc: playerUc}
}

func (h *DashboardHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/me", middleware.RequireAuth(), middleware.AddCacheControlHeader(600, middleware.TYPE_PRIVATE))
	{
		group.GET("/statistics", h.getPlayerStatistics)
		group.GET("/incoming-matches", h.getIncomingMaches)

	}

}

func (h *DashboardHandler) getPlayerStatistics(c *gin.Context) {

	var req struct {
		Limit *int `form:"limit" binding:"omitempty,numeric"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	playerId, exists := c.Get("PlayerId")

	if exists == false || playerId == nil || playerId.(string) == "" {
		res := delivery.NewSuccessResponse(nil)
		c.JSON(http.StatusOK, res)
		return
	}

	statistics, err := h.playerUc.GetPlayerStatistics(c.Request.Context(), player.PlayerStatisticsRequest{
		PlayerId: playerId.(string),
		Limit:    req.Limit,
	})

	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	var response struct {
		EarnedSinglePoints int `json:"earnedSinglePoints"`
		EarnedDoublePoints int `json:"earnedDoublePoints"`
		SinglePoints       int `json:"singlePoints"`
		DoublePoints       int `json:"doublePoints"`
	}

	response.EarnedDoublePoints = statistics.LastDoublePointsSum
	response.EarnedSinglePoints = statistics.LastSinglePointsSum
	response.SinglePoints = statistics.CurrentSinglePoint
	response.DoublePoints = statistics.CurrentDoublePoint

	res := delivery.NewSuccessResponse(response)

	c.JSON(http.StatusOK, res)
}

func (h *DashboardHandler) getIncomingMaches(c *gin.Context) {
	var req struct {
		Limit int16 `form:"limit" binding:"omitempty,numeric"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	playerId, exists := c.Get("PlayerId")

	if exists == false || playerId == nil || playerId.(string) == "" {
		res := delivery.NewSuccessResponse(nil)
		c.JSON(http.StatusOK, res)
		return
	}

	dto := player.PlayerIncomingMatchesRequest{PlayerId: playerId.(string), Limit: req.Limit}
	matches, err := h.playerUc.GetImconimgMatches(c.Request.Context(), dto)

	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	type response struct {
		MatchId      string             `json:"matchId"`
		MatchDate    *time.Time         `json:"matchDate"`
		MatchType    match.Match_TYPE   `json:"matchType"`
		Source       match.Match_SOURCE `json:"source"`
		LeagueId     *string            `json:"leagueId"`
		LeagueName   *string            `json:"leagueName"`
		OppenentId   string             `json:"oppenentId"`
		OppenentName string             `json:"oppenentName"`
	}

	matchesResponse := make([]response, 0, len(matches))
	for _, m := range matches {
		matchesResponse = append(matchesResponse,
			response{MatchId: m.MatchId,
				MatchDate:    m.MatchDate,
				MatchType:    m.MatchType,
				Source:       m.Source,
				LeagueId:     m.LeagueId,
				LeagueName:   m.LeagueName,
				OppenentId:   m.OppenentId,
				OppenentName: m.OppenentName,
			})
	}
	res := delivery.NewSuccessResponse(matchesResponse)

	c.JSON(http.StatusOK, res)

}

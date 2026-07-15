package dashboardhandler

import (
	"tennis-league/common/http/router"
	"tennis-league/user-service/internal/service/player"

	httpcache "tennis-league/common/lib/http/http-cache"
	authmiddleware "tennis-league/common/security/authmiddleware"
)

type DashboardHandler struct {
	playerUc *player.Usecase
}

func NewDashboardHandler(playerUc *player.Usecase) *DashboardHandler {
	return &DashboardHandler{playerUc: playerUc}
}

func (h *DashboardHandler) RegisterRoutes(r *router.CustomRouterGroup) {

	group := r.Group("/me", authmiddleware.RequireAuth(), httpcache.AddCacheControlHeader(600, httpcache.TYPE_PRIVATE))
	{
		group.GET("/statistics", h.getPlayerStatistics)

	}
}

func (h *DashboardHandler) getPlayerStatistics(c *router.CustomContext) {

	var req struct {
		Limit *int `form:"limit" binding:"omitempty,numeric"`
	}

	if !c.BindQueryOrAbort(&req) {
		return
	}

	playerId, exists := c.CurrentPlayerId()

	if exists == false {
		c.OkComplete(nil)
		return
	}

	statistics, err := h.playerUc.GetPlayerStatistics(c.Request.Context(), player.PlayerStatisticsRequest{
		PlayerId: playerId,
		Limit:    req.Limit,
	})

	if err != nil {
		c.ErrorComplete(err)
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
	c.OkComplete(response)
}

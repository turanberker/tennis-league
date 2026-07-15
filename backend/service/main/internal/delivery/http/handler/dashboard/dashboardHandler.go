package dashboard

import (
	"tennis-league/common/http/router"

	"time"

	customerror "tennis-league/common/lib/error"
	httpcache "tennis-league/common/lib/http/http-cache"
	authmiddleware "tennis-league/common/security/authmiddleware"
	"tennis-league/service/internal/domain/match"
)

type DashboardHandler struct {
	matchUseCase *match.UseCase
}

func NewDashboardHandler(matchUseCase *match.UseCase) *DashboardHandler {
	return &DashboardHandler{matchUseCase: matchUseCase}
}

func (h *DashboardHandler) RegisterRoutes(r *router.CustomRouterGroup) {

	group := r.Group("/me", authmiddleware.RequireAuth(), authmiddleware.RequirePlayerRecord, httpcache.AddCacheControlHeader(600, httpcache.TYPE_PRIVATE))
	{
		group.GET("/incoming-matches", h.getIncomingMatches)
	}
}

func (h *DashboardHandler) getIncomingMatches(c *router.CustomContext) {
	var req struct {
		Limit int16 `form:"limit" binding:"omitempty,numeric"`
	}

	if !c.BindQueryOrAbort(&req) {
		return
	}

	playerId, _ := c.CurrentPlayerId()

	dto := match.PlayerIncomingMatchesRequest{PlayerId: playerId, Limit: req.Limit}
	matches, err := h.matchUseCase.GetImconimgMatches(c.Request.Context(), dto)

	if err != nil {
		_ = c.Error(customerror.NewInternalError(err))
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
	c.OkComplete(matchesResponse)
}

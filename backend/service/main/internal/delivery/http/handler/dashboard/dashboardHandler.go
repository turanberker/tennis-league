package dashboard

import (
	"net/http"

	"time"

	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	httpcache "tennis-league/common/lib/http/http-cache"
	authmiddleware "tennis-league/common/security/authmiddleware"
	"tennis-league/service/internal/domain/match"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	matchUsecase *match.UseCase
}

func NewDashboardHandler(matchUsecase *match.UseCase) *DashboardHandler {
	return &DashboardHandler{matchUsecase: matchUsecase}
}

func (h *DashboardHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/me", authmiddleware.RequireAuth(), httpcache.AddCacheControlHeader(600, httpcache.TYPE_PRIVATE))
	{
		group.GET("/incoming-matches", h.getIncomingMaches)

	}
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

	dto := match.PlayerIncomingMatchesRequest{PlayerId: playerId.(string), Limit: req.Limit}
	matches, err := h.matchUsecase.GetImconimgMatches(c.Request.Context(), dto)

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

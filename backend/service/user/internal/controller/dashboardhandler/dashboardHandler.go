package dashboardhandler

import (
	"net/http"
	"tennis-league/user-service/internal/service/player"

	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	httpcache "tennis-league/common/lib/http/http-cache"
	authmiddleware "tennis-league/common/security/authmiddleware"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	playerUc *player.Usecase
}

func NewDashboardHandler(playerUc *player.Usecase) *DashboardHandler {
	return &DashboardHandler{playerUc: playerUc}
}

func (h *DashboardHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/me", authmiddleware.RequireAuth(), httpcache.AddCacheControlHeader(600, httpcache.TYPE_PRIVATE))
	{
		group.GET("/statistics", h.getPlayerStatistics)

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

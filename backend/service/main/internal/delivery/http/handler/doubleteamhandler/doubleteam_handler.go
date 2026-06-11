package doubleteamhandler

import (
	"net/http"

	"tennis-league/common/lib/http/delivery"
	httpcache "tennis-league/common/lib/http/http-cache"
	"tennis-league/service/internal/delivery/http/handler/playerhandler"
	"tennis-league/service/internal/domain/team"

	"github.com/gin-gonic/gin"
)

type DoubleTeamHandler struct {
	uc *team.UseCase
}

func NewDoubleTeamHandler(uc *team.UseCase) *DoubleTeamHandler {
	return &DoubleTeamHandler{uc: uc}
}

func (h *DoubleTeamHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/double-team")
	{
		group.GET("/:id/members", httpcache.AddCacheControlHeader(600, httpcache.TYPE_PUBLIC), h.getTeamMembers)
	}

}

func (h *DoubleTeamHandler) getTeamMembers(c *gin.Context) {
	id := c.Param("id")

	players, err := h.uc.GetTeamMembers(c.Request.Context(), id)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var response []playerhandler.PlayerResponse
	for _, p := range players {
		response = append(response, playerhandler.PlayerResponse{
			ID:           p.ID,
			Name:         p.Name,
			Surname:      p.Surname,
			Sex:          p.Sex,
			UserId:       p.UserId,
			DoublePoints: p.DoublePoints,
			SinglePoints: p.SinglePoints,
		})
	}

	res := delivery.NewSuccessResponse(response)
	c.JSON(http.StatusOK, res)
}

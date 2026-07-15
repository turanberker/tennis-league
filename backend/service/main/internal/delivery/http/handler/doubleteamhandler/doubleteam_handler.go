package doubleteamhandler

import (
	"tennis-league/common/http/router"

	httpcache "tennis-league/common/lib/http/http-cache"
	"tennis-league/service/internal/domain/team"
)

type DoubleTeamHandler struct {
	uc *team.UseCase
}

func NewDoubleTeamHandler(uc *team.UseCase) *DoubleTeamHandler {
	return &DoubleTeamHandler{uc: uc}
}

func (h *DoubleTeamHandler) RegisterRoutes(r *router.CustomRouterGroup) {

	group := r.Group("/double-team")
	{
		group.GET("/:id/members", httpcache.AddCacheControlHeader(600, httpcache.TYPE_PUBLIC), h.getTeamMembers)
	}

}

func (h *DoubleTeamHandler) getTeamMembers(c *router.CustomContext) {
	id := c.Param("id")

	players, err := h.uc.GetTeamMembers(c.Request.Context(), id)

	if err != nil {
		c.ErrorComplete(err)
		return
	}

	var response []PlayerResponse
	for _, p := range players {
		response = append(response, PlayerResponse{
			ID:           p.ID,
			Name:         p.Name,
			Surname:      p.Surname,
			Sex:          p.Sex,
			UserId:       p.UserId,
			DoublePoints: p.DoublePoints,
			SinglePoints: p.SinglePoints,
		})
	}
	c.OkComplete(response)
}

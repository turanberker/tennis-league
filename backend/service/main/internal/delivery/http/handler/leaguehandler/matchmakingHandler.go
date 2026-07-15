package leaguehandler

import (
	"tennis-league/common/http/router"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/service/internal/domain/matchrequest"
)

type matchmakingHandler struct {
	leagueHandlerMiddleware *leagueHandlerMiddleware
	useCase                 *matchrequest.UseCase
}

func newMatchMakingHandler(leagueHandlerMiddleware *leagueHandlerMiddleware,
	useCase *matchrequest.UseCase) *matchmakingHandler {
	return &matchmakingHandler{leagueHandlerMiddleware: leagueHandlerMiddleware, useCase: useCase}
}

func (h *matchmakingHandler) registerSubRoutes(group *router.CustomRouterGroup) {
	group.POST("", authmiddleware.RequireAuth(), authmiddleware.RequirePlayerRecord,
		h.leagueHandlerMiddleware.userAttendedToLeague, h.new)
}

func (h *matchmakingHandler) new(c *router.CustomContext) {
	leagueId := c.Param("id")
	playerId, _ := c.CurrentPlayerId()
	request := MatchRequestRequest{}

	if !c.BindJSONOrAbort(&request) {
		return
	}

	dto := matchrequest.RequestDto{
		LeagueId:  &leagueId,
		PlayerId:  playerId,
		Date:      request.MatchRequestDate,
		StartHour: request.StartHour,
	}

	requestId, err := h.useCase.NewRequest(c.Request.Context(), dto)
	if err != nil {
		c.ErrorComplete(err)
		return
	}
	c.OkComplete(requestId)
}

package leaguehandler

import (
	"errors"
	"net/http"
	"tennis-league/common/http/router"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/service/internal/domain/matchrequest"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (h *matchmakingHandler) new(c *gin.Context) {
	leagueId := c.Param("id")
	playerId, _ := authmiddleware.GetPlayerIdFromContext(c)
	request := MatchRequestRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			_ = c.Error(customerror.NewValidationError(ve))
			c.Abort()
			return
		}
	}

	dto := matchrequest.RequestDto{
		LeagueId:  &leagueId,
		PlayerId:  playerId,
		Date:      request.MatchRequestDate,
		StartHour: request.StartHour,
	}

	requestId, err := h.useCase.NewRequest(c.Request.Context(), dto)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, delivery.NewSuccessResponse(requestId))
}

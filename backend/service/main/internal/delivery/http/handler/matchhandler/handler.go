package matchhandler

import (
	"errors"
	"log"
	"net/http"
	"time"

	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"
	authmiddleware "tennis-league/common/security/auth"

	errorcodes "tennis-league/service/internal/domain/error_codes"
	"tennis-league/service/internal/domain/match"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MatchHandler struct {
	u *match.UseCase
}

func NewMatchHandler(u *match.UseCase) *MatchHandler {
	return &MatchHandler{u: u}
}

func (h *MatchHandler) RegisterRoutes(r *gin.Engine) {
	matches := r.Group("/match")
	{
		matches.GET("/:id", h.getById)
		matches.GET("/:id/match-info", h.getSetScore)
		matches.PUT("/:id/score", authmiddleware.RequireAuth(), h.checkIfMatchIsFriendly, h.checkIfUserIsMatchPlayer, h.updateScore)
		matches.PUT("/:id/update-date", authmiddleware.RequireAuth(), h.checkIfMatchIsFriendly, h.checkIfUserIsMatchPlayer, h.updateDate)
		matches.PUT("/:id/approve", authmiddleware.RequireAuth(), h.checkIfMatchIsFriendly, h.checkIfUserIsMatchPlayer, h.approveScore)
	}
}
func (h *MatchHandler) approveScore(c *gin.Context) {
	matchId := c.Param("id")
	err := h.u.ApproveScore(c.Request.Context(), match.MatchSource_FRIENDLY, matchId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(nil))
}
func (h *MatchHandler) getSetScore(c *gin.Context) {

	// path param
	matchId := c.Param("id")
	setScores := h.u.GetSetScore(c.Request.Context(), matchId)
	sides, err := h.u.GetMatchInfo(c.Request.Context(), matchId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	response := MatchSetScoreResponse{}

	response.MatchInfo.MatchDate = sides.MatchDate
	response.MatchInfo.MatchType = sides.MatchType
	response.MatchInfo.Source = sides.Source
	response.MatchInfo.SourceId = sides.SourceId
	response.MatchInfo.Status = sides.Status
	response.MatchInfo.Side1 = sides.Side1.Name
	response.MatchInfo.Side2 = sides.Side2.Name
	for _, s := range setScores {
		switch s.SetNumber {
		case 1:
			if s.Team1Game != nil {
				response.MatchScore.Set1.Team1Score = *s.Team1Game
			}
			if s.Team2Game != nil {
				response.MatchScore.Set1.Team2Score = *s.Team2Game
			}
		case 2:
			if s.Team1Game != nil {
				response.MatchScore.Set2.Team1Score = *s.Team1Game
			}
			if s.Team2Game != nil {
				response.MatchScore.Set2.Team2Score = *s.Team2Game
			}

		case 3:
			response.MatchScore.SuperTie = &SetScore{}
			if s.Team1TiePoint != nil {
				response.MatchScore.SuperTie.Team1Score = *s.Team1TiePoint
			}
			if s.Team2TiePoint != nil {
				response.MatchScore.SuperTie.Team2Score = *s.Team2TiePoint
			}
		}

	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(response))
}

func (h *MatchHandler) getById(c *gin.Context) {

	// path param
	matchId := c.Param("id")
	log.Printf("match id: %s", matchId)
	c.JSON(200, gin.H{"message": "get match by id"})
}

func (h *MatchHandler) updateScore(c *gin.Context) {
	matchId := c.Param("id")

	macScore := UpdateScoreRequest{}

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

	response, err := h.u.SaveMatchScore(c.Request.Context(), saveMatchScore)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(&MatchScoreResponse{Team1Score: response.Team1Score, Team2Score: response.Team2Score}))

}

func (h *MatchHandler) updateDate(c *gin.Context) {
	matchId := c.Param("id")

	var req struct {
		MatchDate time.Time `form:"match-date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
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

	err := h.u.UpdateMatchDate(c.Request.Context(), matchId, match.MatchSource_FRIENDLY, &req.MatchDate)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(req.MatchDate))

}

func (h *MatchHandler) checkIfMatchIsFriendly(c *gin.Context) {
	matchId := c.Param("id")
	matchInfo, err := h.u.GetMatchInfo(c.Request.Context(), matchId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	if matchInfo.Source != match.MatchSource_FRIENDLY {
		c.Error(errors.New("Buradan sadece dosluk maçları güncellenebilir"))
		c.Abort()
	}
}

func (h *MatchHandler) checkIfUserIsMatchPlayer(c *gin.Context) {
	matchId := c.Param("id")
	playerId, exists := authmiddleware.GetPlayerIdFromContext(c)
	if !exists {
		err := &customerror.BusinnesException{
			StatusCode: http.StatusForbidden,
			ErrorCode:  errorcodes.INSUFFICIENT_PERMISSIONS,
			Message:    "Oyuncu kaydınız bulunamamıştır",
		}
		c.Error(err)
		c.Abort()
	}

	playedInMatch, err := h.u.IsUserPlayerOfMatch(c.Request.Context(), matchId, playerId)

	if err != nil {
		c.Error(err)
		c.Abort()
	}

	if !playedInMatch {
		err := &customerror.BusinnesException{
			StatusCode: http.StatusForbidden,
			ErrorCode:  errorcodes.ErrNotParticipatedInMatch,
			Message:    "Bu maçta oynamadığınız için skoru güncelleyemezsiniz",
		}
		c.Error(err)
		c.Abort()
	}
}

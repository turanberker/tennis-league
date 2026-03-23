package matchhandler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/turanberker/tennis-league-service/internal/delivery"
	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
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
		matches.GET("/:id/set-scores", h.getSetScore)
		matches.PUT("/:id/score", h.updateScore)
		matches.PUT("/:id/update-date", h.updateDate)
		matches.PUT("/:id/approve", h.approveScore)
	}
}
func (h *MatchHandler) approveScore(c *gin.Context) {
	matchId := c.Param("id")
	err := h.u.ApproveScore(c.Request.Context(), match.MatchSource_FRIENDLY, matchId)
	if err != nil {
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(nil))
}
func (h *MatchHandler) getSetScore(c *gin.Context) {

	// path param
	matchId := c.Param("id")
	setScores := h.u.GetSetScore(c.Request.Context(), matchId)

	response := MatchScore{}
	for _, s := range setScores {
		switch s.SetNumber {
		case 1:
			if s.Team1Game != nil {
				response.Set1.Team1Score = *s.Team1Game
			}
			if s.Team2Game != nil {
				response.Set1.Team2Score = *s.Team2Game
			}
		case 2:
			if s.Team1Game != nil {
				response.Set2.Team1Score = *s.Team1Game
			}
			if s.Team2Game != nil {
				response.Set2.Team2Score = *s.Team2Game
			}

		case 3:
			response.SuperTie = &SetScore{}
			if s.Team1TiePoint != nil {
				response.SuperTie.Team1Score = *s.Team1TiePoint
			}
			if s.Team2TiePoint != nil {
				response.SuperTie.Team2Score = *s.Team2TiePoint
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

	macScore := MatchScore{}

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

	saveMatchScore := &match.SaveMatchScore{MatchId: matchId, Set1: set1, Set2: set2}

	if macScore.SuperTie != nil {

		saveMatchScore.SuperTie = &match.SaveScore{}
		saveMatchScore.SuperTie.Team1Score = macScore.SuperTie.Team1Score
		saveMatchScore.SuperTie.Team2Score = macScore.SuperTie.Team2Score

	}

	response, err := h.u.SaveMatchScore(c.Request.Context(), saveMatchScore)
	if err != nil {
		c.Error(customerror.NewInternalError(err))
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
		c.Error(customerror.NewInternalError(err))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(req.MatchDate))

}

package matchhandler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/turanberker/tennis-league-service/internal/delivery"
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
		matches.PUT("/:id/score", h.updateScore)
		matches.PUT("/:id/update-date", h.updateDate)
	}
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
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		return
	}

	log.Printf("match id: %s", matchId)
	log.Printf("score :%+v", macScore)

	set1 := match.SaveScore{Team1Score: macScore.Set1.Team1Score, Team2Score: macScore.Set1.Team2Score}
	set2 := match.SaveScore{Team1Score: macScore.Set2.Team1Score, Team2Score: macScore.Set2.Team2Score}

	var tie match.SaveScore
	if macScore.SuperTie != nil {
		tie.Team1Score = macScore.SuperTie.Team1Score
		tie.Team2Score = macScore.SuperTie.Team1Score
	}

	response, err := h.u.SaveMatchScore(c.Request.Context(), &match.SaveMatchScore{MatchId: matchId, Set1: set1, Set2: set2, SuperTie: &tie})
	if err != nil {
		c.JSON(http.StatusBadRequest, delivery.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery.NewSuccessResponse(&MatchScoreResponse{Team1Score: response.Team1Score, Team2Score: response.Team2Score}))

}

func (h *MatchHandler) updateDate(c *gin.Context) {
	matchId := c.Param("id")
	matchDateString := c.Query("match-date")
	var matchDate *time.Time
	if matchDateString != "" {
		t, err := time.Parse(time.RFC3339, matchDateString)
		if err != nil {
			c.JSON(http.StatusBadRequest, delivery.NewErrorResponse("Tarih Formatı Hatalı"))
			return
		}
		matchDate = &t
	}

	h.u.UpdateMatchDate(c.Request.Context(), matchId, matchDate)
	log.Printf("match id: %s, match Date: %s", matchId, matchDate)
	c.JSON(http.StatusOK, delivery.NewSuccessResponse(nil))

}

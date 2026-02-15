package matchhandler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{}
}

func (h *MatchHandler) RegisterRoutes(r *gin.Engine) {
	matches := r.Group("/match")
	{
		matches.GET("/:id", h.getById)
		matches.PATCH("/:id", h.updateScore)
		matches.PATCH("/:id/update-date", h.updateDate)
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
	log.Printf("match id: %s", matchId)
	c.JSON(200, gin.H{"message": "get match by id"})
	// path param
}

func (h *MatchHandler) updateDate(c *gin.Context) {
	matchId := c.Param("id")
	matchDate := c.Query("match-date")
	log.Printf("match id: %s, match Date: %s", matchId, matchDate)
	c.JSON(200, gin.H{"message": "get match by id"})
	
}

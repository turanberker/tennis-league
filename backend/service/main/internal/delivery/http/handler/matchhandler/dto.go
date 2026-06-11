package matchhandler

import (
	"time"

	"tennis-league/service/internal/domain/match"
)

type SetScore struct {
	Team1Score int8 `json:"team1Score" binding:"gte=0,lte=99"`
	Team2Score int8 `json:"team2Score" binding:"gte=0,lte=99"`
}

type UpdateScoreRequest struct {
	MatchDate *time.Time `json:"matchDate" binding:"required"`
	MatchScore
}

type MatchSetScoreResponse struct {
	MatchInfo  MatchInfo  `json:"matchInfo"`
	MatchScore MatchScore `json:"setScore"`
}

type MatchScore struct {
	Set1     SetScore  `json:"set1" binding:"required,tennis_set"`
	Set2     SetScore  `json:"set2" binding:"required,tennis_set"`
	SuperTie *SetScore `json:"superTie" binding:"omitempty,super_tie"`
}

type MatchScoreResponse struct {
	Team1Score int8 `json:"team1Score"`
	Team2Score int8 `json:"team2Score"`
}

type MatchInfo struct {
	MatchDate *time.Time         `json:"matchDate"`
	Source    match.Match_SOURCE `json:"source"`
	SourceId  *string            `json:"sourceId"` //leagueId veya TournamentId olacak
	MatchType match.Match_TYPE   `json:"type"`
	Status    match.MATCH_Status `json:"status"`
	Side1     string             `json:"side1"`
	Side2     string             `json:"side2"`
}

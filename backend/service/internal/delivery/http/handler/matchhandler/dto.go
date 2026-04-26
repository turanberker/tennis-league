package matchhandler

import "time"

type SetScore struct {
	Team1Score int8 `json:"team1Score" binding:"gte=0,lte=99"`
	Team2Score int8 `json:"team2Score" binding:"gte=0,lte=99"`
}

type UpdateScoreRequest struct {
	MatchDate *time.Time `json:"matchDate" binding:"omitempty"`
	MatchScore
}

type MatchSetScoreResponse struct {
	MatchDate *time.Time `json:"matchDate"`
	MatchScore
	Side1 string `json:"side1"`
	Side2 string `json:"side2"`
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

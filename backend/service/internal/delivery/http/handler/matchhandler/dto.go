package matchhandler

type SetScore struct {
	Team1Score int8 `json:"team1Score" binding:"gte=0,lte=99"`
	Team2Score int8 `json:"team2Score" binding:"gte=0,lte=99"`
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

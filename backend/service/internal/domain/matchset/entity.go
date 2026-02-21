package matchSet

type UpdateSetScore struct {
	MatchId    string
	Set        int8
	Team1Score int8
	Team2Score int8
}

type UpdateSuperTieScore struct {
	MatchId    string
	Team1Score int8
	Team2Score int8
}

type MatchSetScores struct {
	SetNumber     int8
	Team1Game     *int8
	Team2Game     *int8
	Team1TiePoint *int8
	Team2TiePoint *int8
}

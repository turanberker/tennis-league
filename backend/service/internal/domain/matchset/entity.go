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

package scoreboard

type TeamReferance struct {
	Id   string
	Name string
}

type ScoreBoard struct {
	Team      TeamReferance
	Played    int16
	Won       int16
	Lost      int16
	WonSets   int16
	LostSets  int16
	WonGames  int16
	LostGames int16
	Score     int16
}

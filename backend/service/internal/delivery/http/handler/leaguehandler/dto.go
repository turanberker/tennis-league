package leaguehandler

import (
	"time"

	"github.com/turanberker/tennis-league-service/internal/domain/match"
)

type LeagueResponse struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	FixtureCreatedDate *time.Time `json:"fixtureCreatedDate,omitempty"`
}

type LeagueFixtureMatchResponse struct {
	Id        string          `json:"id"`
	Team1     TeamRefResponse `json:"team1"`
	Team2     TeamRefResponse `json:"team2"`
	Status    match.Status    `json:"status"`
	MatchDate *time.Time      `json:"matchDate,omitempty"`
}

type TeamRefResponse struct {
	TeamRef
	Score  *int8 `json:"score"`
	Winner *bool `json:"winner"`
}

type ScoreBoardResponse struct {
	TeamRef
	Order     int `json:"order"`
	Played    int16 `json:"played"`
	Won       int16 `json:"won"`
	Lost      int16 `json:"lost"`
	WonSets   int16 `json:"wonSets"`
	LostSets  int16 `json:"lostSets"`
	WonGames  int16 `json:"wonGames"`
	LostGames int16 `json:"lostGames"`
	Score     int16 `json:"score"`
}

type TeamRef struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

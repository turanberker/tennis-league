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
	Id   string `json:"id"`
	Name string `json:"name"`
}

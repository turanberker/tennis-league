package match

import (
	"time"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusApproved  Status = "SCORE_APPROVED"
	StatusCancelled Status = "CANCELLED"
)

type PersistLeagueMatch struct {
	LeagueId string
	Team1Id  string
	Team2Id  string
}

type UpdateMatchDate struct {
	Id        string
	MatchDate *time.Time
}

type MatchTeamIds struct {
	Team1Id string
	Team2Id string
	Status  Status
}

type UpdateMatchScore struct {
	Id           string
	Team1Score   int8
	Team2Score   int8
	WinnerTeamId string
}

type LeagueFixtureMatch struct {
	Id        string
	LeagueId  string
	Team1     teamRef
	Team2     teamRef
	Status    Status
	MatchDate *time.Time
}

type teamRef struct {
	Id     string
	Name   string
	Score  *int8
	Winner *bool
}

type MatchApprovedEvent struct {
	MatchID string `json:"matchId"`
}

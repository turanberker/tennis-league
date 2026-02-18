package match

import (
	"time"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	RoleCancelled   Status = "CANCELLED"
)

type PersistLeagueMatch struct {
	LeagueId string
	Team1Id  string
	Team2Id  string
}

type UpdateMatchDate struct{
	Id string
	MatchDate *time.Time
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
	Id   string
	Name string
}

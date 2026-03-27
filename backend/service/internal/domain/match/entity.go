package match

import (
	"time"
)

type Status string
type Match_TYPE string
type Match_SOURCE string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusApproved  Status = "SCORE_APPROVED"
	StatusCancelled Status = "CANCELLED"

	MatchType_SINGLE Match_TYPE = "SINGLE"
	MatchType_DOUBLE Match_TYPE = "DOUBLE"

	MatchSource_FRIENDLY   Match_SOURCE = "FRIENDLY"
	MatchSource_TOURNAMENT Match_SOURCE = "TOURNAMENT"
	MatchSource_PLAYOFF    Match_SOURCE = "PLAYOFF"
)

type BulkInsertMatches struct {
	Sides []SideIds
	Type  MatchType
}

type SideIds struct {
	Side1 string
	Side2 string
}

type MatchType struct {
	Id     *string //LEAGUEID veya TournamentId
	Source Match_SOURCE
	Type   Match_TYPE
}

type UpdateMatchDate struct {
	Id        string
	Source    Match_SOURCE
	MatchDate *time.Time
}

type MatchTeamIds struct {
	LeagueId string
	Team1Id  string
	Team2Id  string
	Status   Status
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

type MatchParticipant struct {
	PlayerID    string
	DoublePoint int
	IsWinner    bool
}

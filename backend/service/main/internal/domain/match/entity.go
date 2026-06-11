package match

import (
	"time"
)

type MATCH_Status string
type Match_TYPE string
type Match_SOURCE string

const (
	StatusPending   MATCH_Status = "PENDING"
	StatusCompleted MATCH_Status = "COMPLETED"
	StatusApproved  MATCH_Status = "SCORE_APPROVED"
	StatusCancelled MATCH_Status = "CANCELLED"

	MatchType_SINGLE Match_TYPE = "SINGLE"
	MatchType_DOUBLE Match_TYPE = "DOUBLE"
	MatchType_TEAM   Match_TYPE = "TEAM"

	MatchSource_FRIENDLY Match_SOURCE = "FRIENDLY"
	MatchSource_LEAGUE   Match_SOURCE = "LEAGUE"
	MatchSource_PLAYOFF  Match_SOURCE = "PLAYOFF"
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
	MatchDate *time.Time
}

type MatchTeamIds struct {
	LeagueId string
	Team1Id  string
	Team2Id  string
	Status   MATCH_Status
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
	Team1     TeamRef
	Team2     TeamRef
	Status    MATCH_Status
	MatchDate *time.Time
}

type MatchInfo struct {
	MatchDate *time.Time
	Side1     MatchSide
	Side2     MatchSide
	Source    Match_SOURCE
	SourceId  *string
	MatchType Match_TYPE
	Status    MATCH_Status
}

type MatchSide struct {
	Id   string
	Name string
}

type TeamRef struct {
	MatchSide
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

type PlayerIncomingMatchesQueryParam struct {
	PlayerId string
	Limit    int16
}

type PlayerIncomingMatchesResult struct {
	MatchId      string
	MatchDate    *time.Time
	MatchType    Match_TYPE
	Source       Match_SOURCE
	LeagueId     *string
	LeagueName   *string
	OppenentId   string
	OppenentName string
}

type FixtureFilter struct {
	TeamId *string
}

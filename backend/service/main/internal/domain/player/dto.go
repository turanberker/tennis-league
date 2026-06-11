package player

import (
	"time"

	"tennis-league/service/internal/domain/match"
)

type PlayerIncomingMatchesRequest struct {
	PlayerId string
	Limit    int16
}

type IncomingMatches struct {
	MatchId      string
	MatchDate    *time.Time
	MatchType    match.Match_TYPE
	Source       match.Match_SOURCE
	LeagueId     *string
	LeagueName   *string
	OppenentId   string
	OppenentName string
}

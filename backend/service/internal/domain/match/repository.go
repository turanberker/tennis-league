package match

import (
	"context"
)

type Repository interface {
	SaveLeagueMatches(ctx context.Context, matches []*PersistLeagueMatch) error
	GetFixtureByLeagueId(ctx context.Context, leagueId string) ([]*LeagueFixtureMatch, error)
	UpdateMatchDate(ctx context.Context, data UpdateMatchDate) error
	GetMatchTeamIds(ctx context.Context, matchId string) *MatchTeamIds
	UpdateMatchScore(ctx context.Context, macScore *UpdateMatchScore) error
	ApproveScore(ctx context.Context, matchId string) error
}

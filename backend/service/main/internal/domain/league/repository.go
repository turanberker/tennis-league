package league

import (
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id string) (*League, error)

	GetAll(ctx context.Context, status *LEAGUE_STATUS) ([]*LeagueListSelect, error)

	Save(ctx context.Context, persistLeague *PersistLeague) (*string, error)

	StartLeague(ctx context.Context, leagueId string) error

	IsFixtureCreated(ctx context.Context, leagueId string) (bool, error)

	IncreaseAttendanceCount(ctx context.Context, leagueId string) (*int32, error)
}

type ParticipantRepository interface {
	SingleLeagueAttendanceList(ctx context.Context, leagueId string) ([]SingleLeagueAttendance, error)

	AddPlayerToLeague(ctx context.Context, leagueAttendance NewLeaguePlayerAttendance) error

	AddTeamToLeague(ctx context.Context, leagueId string, teamId string) error

	IsPlayerAttendedToLeague(ctx context.Context, leagueId string, playerId string) (bool, error)
}

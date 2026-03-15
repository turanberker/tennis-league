package leaguecoordinator

import "context"

type Repository interface {
	Exists(ctx context.Context, leagueId string, userId string) (bool, error)
	Add(ctx context.Context, leagueId string, userId string) (*bool, error)
}

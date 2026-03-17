package player

import (
	"context"
)

type Repository interface {
	GetById(ctx context.Context, id int64) (*Player, error)

	Save(ctx context.Context, persistPlayer *PersistPlayer) (*string, error)

	List(ctx context.Context, queryParams ListQueryParameters) ([]*Player, error)

	AssignToUser(ctx context.Context, playerId string, userId string) error
}

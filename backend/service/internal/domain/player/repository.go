package player

import "context"

type Repository interface {
	GetById(ctx context.Context, id int64) (*Player, error)

	GetByUuid(ctx context.Context, uuid string) (*Player, error)

	Save(ctx context.Context, persistPlayer *PersistPlayer) (int64, error)

	List(ctx context.Context, name string) ([]*Player, error)	
}

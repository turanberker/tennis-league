package league

import "context"

type Repository interface {
	GetById(ctx context.Context, id int64) (*League, error)

	GetAll(ctx context.Context, name string) ([]*League, error)

	Save(ctx context.Context, persistLeague *PersistLeague) (int64, error)
}

package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*LoginUserCheck, error)

	ExistsByEmail(ctx context.Context, email string) bool

	SaveUser(ctx context.Context, tx *sql.Tx, u *PersistUser) (string, error)

	List(ctx context.Context) ([]*User, error)

	UpdateRoleAsCoordinator(txCtx context.Context, userId string) error
}

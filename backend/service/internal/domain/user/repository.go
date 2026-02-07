package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)

	ExistsByEmail(ctx context.Context, email string) bool

	SaveUser(ctx context.Context, tx *sql.Tx, u *User) (*User, error)
}

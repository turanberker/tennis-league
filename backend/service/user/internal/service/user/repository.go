package user

import (
	"context"
)

type Repository interface {
	FindById(ctx context.Context, id string) (*UserData, error)

	GetByEmail(ctx context.Context, email string) (*LoginUserCheck, error)

	ExistsByEmail(ctx context.Context, email string) bool

	SaveUser(ctx context.Context, u *PersistUser) (string, error)

	List(ctx context.Context) ([]*User, error)

	UpdateRoleAsCoordinator(txCtx context.Context, userId string) error

	UpdatePassword(ctx context.Context, userId string, newPasswordHash string) error
}

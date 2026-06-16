package user

import (
	"context"
)

type Repository interface {
	UpdateRoleAsCoordinator(txCtx context.Context, userId string) error
}

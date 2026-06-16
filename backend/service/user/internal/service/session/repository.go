package session

import (
	"context"
	"tennis-league/common/security/dto"
)

type Repository interface {
	Start(ctx context.Context, startSessionInput *StartSessionInput) (*dto.Session, error)

	Get(ctx context.Context, sessionID string) (*dto.Session, error)

	Delete(ctx context.Context, sessionID string) error

	Refresh(ctx context.Context, sessionID string) error
}

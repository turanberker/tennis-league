package session

import "context"

type Repository interface {
	Start(ctx context.Context, userID int64, role string) (*Session, error)

	Get(ctx context.Context, sessionID string) (*Session, error)

	Delete(ctx context.Context, sessionID string) error

	Refresh(ctx context.Context, sessionID string) error
}
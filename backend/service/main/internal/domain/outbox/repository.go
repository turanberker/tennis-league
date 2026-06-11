package outbox

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, entity *PersistEntity) error

	GetEventForUpdate(ctx context.Context, id string) (*EventToPublish, error)

	GetPendingIDs(ctx context.Context, limit int) ([]string, error)

	IncreaseRetryCount(ctx context.Context, id string) error

	UpdateToPublished(ctx context.Context, id string) error
}

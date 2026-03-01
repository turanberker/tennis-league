package outbox

import (
	"context"
	"database/sql"
)

type Repository interface {
	Save(ctx context.Context, tx *sql.Tx, entity *PersistEntity) error

	GetEventsToPublish(ctx context.Context, tx *sql.Tx) ([]*EventToPublish, error)

	IncreaseRetryCount(ctx context.Context, tx *sql.Tx, id string)

	UpdateToPublished(ctx context.Context, tx *sql.Tx, id string) error
}

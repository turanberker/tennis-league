package postgres

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/domain/outbox"
)

type Repository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, tx *sql.Tx, entity *outbox.PersistEntity) error {

	query := `
		INSERT INTO outbox_events
		( aggregate_type, aggregate_id, event_type, payload, status, retry_count, created_at)
		VALUES ($1,$2,$3,$4,$5,0,current_date)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		entity.AggregateType,
		entity.AggregateID,
		entity.EventType,
		entity.Payload,
		outbox.StatusPending,
	)

	return err
}

func (r *Repository) GetEventsToPublish(ctx context.Context, tx *sql.Tx) ([]*outbox.EventToPublish, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, event_type, payload
		FROM outbox_events
		WHERE status = $1
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT 10
	`, outbox.StatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type event struct {
		id        string
		eventType string
		payload   []byte
	}

	var events []*outbox.EventToPublish

	for rows.Next() {
		var e = &outbox.EventToPublish{}
		if err := rows.Scan(&e.Id, &e.EventType, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *Repository) IncreaseRetryCount(ctx context.Context, tx *sql.Tx, id string) {
	_, _ = tx.ExecContext(ctx,
		`UPDATE outbox_events
				 SET retry_count = retry_count + 1
				 WHERE id = $1`, id)
}

func (r *Repository) UpdateToPublished(ctx context.Context, tx *sql.Tx, id string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE outbox_events
			 SET status = $1,
			     processed_at = current_date
			 WHERE id = $2`, outbox.StatusPublished, id)
	return err
}

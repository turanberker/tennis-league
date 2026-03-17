package postgres

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/domain/outbox"
)

type Repository struct {
	BaseRepository
}

func NewOutboxRepository(db *sql.DB) *Repository {
	return &Repository{BaseRepository{db: db}}
}

func (r *Repository) Save(ctx context.Context, entity *outbox.PersistEntity) error {
	exec := r.GetExecutor(ctx)
	query := `
		INSERT INTO outbox_events
		( aggregate_type, aggregate_id, event_type, payload, status, retry_count, created_at)
		VALUES ($1,$2,$3,$4,$5,0,NOW())
	`

	_, err := exec.ExecContext(
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

func (r *Repository) GetPendingIDs(ctx context.Context, limit int) ([]string, error) {
	exec := r.GetExecutor(ctx)

	query := `
        SELECT id 
        FROM outbox_events 
        WHERE status = $1 
        ORDER BY created_at ASC 
        LIMIT $2`

	rows, err := exec.QueryContext(ctx, query, outbox.StatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) IncreaseRetryCount(ctx context.Context, id string) error {
	exec := r.GetExecutor(ctx)
	query := `
        UPDATE outbox_events 
        SET retry_count = retry_count + 1,
            status = CASE WHEN retry_count + 1 >= 5 THEN 'FAILED' ELSE status END
        WHERE id = $1`

	_, err := exec.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) UpdateToPublished(ctx context.Context, id string) error {
	exec := r.GetExecutor(ctx)
	_, err := exec.ExecContext(ctx,
		`UPDATE outbox_events
			 SET status = $1,
			     processed_at = current_date
			 WHERE id = $2`, outbox.StatusPublished, id)
	return err
}

func (r *Repository) GetEventForUpdate(ctx context.Context, id string) (*outbox.EventToPublish, error) {
	exec := r.GetExecutor(ctx)
	// Sadece bu ID'yi kilitliyoruz
	row := exec.QueryRowContext(ctx, `
        SELECT id, event_type, payload 
        FROM outbox_events 
        WHERE id = $1 FOR UPDATE SKIP LOCKED`, id)

	var e outbox.EventToPublish
	if err := row.Scan(&e.Id, &e.EventType, &e.Payload); err != nil {
		return nil, err
	}
	return &e, nil
}

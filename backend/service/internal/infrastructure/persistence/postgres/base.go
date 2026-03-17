package postgres

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type BaseRepository struct {
	db *sql.DB
}

// GetExecutor, context'e bakar: Transaction varsa onu, yoksa ana DB'yi döner.
// Dikkat: Dönüş tipi QueryExecutor interface'idir!
func (r *BaseRepository) GetExecutor(ctx context.Context) database.QueryExecutor {
	if tx, ok := database.GetTxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

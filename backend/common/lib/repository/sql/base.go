package sqlrepository

import (
	"context"
	"database/sql"
	"tennis-league/common/lib/database"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// GetExecutor, context'e bakar: Transaction varsa onu, yoksa ana DB'yi döner.
// Dikkat: Dönüş tipi QueryExecutor interface'idir!
func (r *Repository) GetExecutor(ctx context.Context) database.QueryExecutor {
	if tx, ok := database.GetTxFromContext(ctx); ok {
		return tx
	}
	return r.DB
}

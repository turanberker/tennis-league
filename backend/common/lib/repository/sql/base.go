package sqlrepository

import (
	"context"
	"database/sql"
	"log"
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

	var exec database.QueryExecutor = r.DB
	if tx, ok := database.GetTxFromContext(ctx); ok {
		return tx
	}
	return LoggingExecutor{QueryExecutor: exec}
}

// --- LOGGING SARMALAYICI (DECORATOR) ---
type LoggingExecutor struct {
	database.QueryExecutor // İçerideki gerçek sql.DB veya sql.Tx nesnesi
}

func (l LoggingExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Printf("[SQL EXEC] Query: %s | Args: %v", query, args)
	return l.QueryExecutor.ExecContext(ctx, query, args...)
}

func (l LoggingExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Printf("[SQL QUERY] Query: %s | Args: %v", query, args)
	return l.QueryExecutor.QueryContext(ctx, query, args...)
}

func (l LoggingExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Printf("[SQL ROW] Query: %s | Args: %v", query, args)
	return l.QueryExecutor.QueryRowContext(ctx, query, args...)
}

package database

import (
	"context"
	"database/sql"
)

// QueryExecutor, 'sorgu çalıştırabilen her şey' demektir.
type QueryExecutor interface {
	// Bu üç metod hem sql.DB'de hem de sql.Tx'de birebir aynı imza ile var.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

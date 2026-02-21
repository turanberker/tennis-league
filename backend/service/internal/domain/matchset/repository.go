package matchSet

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveSetScore(ctx context.Context, tx *sql.Tx, setScore *UpdateSetScore) error
	SaveSuperTieScore(ctx context.Context, tx *sql.Tx, setScore *UpdateSuperTieScore) error
	DeleteSetScores(ctx context.Context, tx *sql.Tx, matchId string) error
}

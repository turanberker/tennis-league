package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/service/internal/delivery/message/consumer/match_score/matchScoreApproved"

	"github.com/Masterminds/squirrel"
)

type PlayerEarnedPointRepository struct {
	sqlrepository.Repository
}

func NewPlayerEarnedPointRepository(db *sql.DB) *PlayerEarnedPointRepository {
	return &PlayerEarnedPointRepository{Repository: *sqlrepository.NewRepository(db)} // Base'deki db'yi dolduruyoruz
}

func (r *PlayerEarnedPointRepository) AddPlayerPoint(ctx context.Context, addPlayerPoint *matchScoreApproved.AddPlayerPoint) error {
	exec := r.GetExecutor(ctx)

	builder := squirrel.Insert("tennisleague.earned_points").
		PlaceholderFormat(squirrel.Dollar).
		Columns("player_id", "match_date", "earned_point", "match_type").
		Values(
			addPlayerPoint.PlayerId,
			time.Now(), // Veya addPlayerPoint içinden gelen tarih
			addPlayerPoint.EarnedPoint,
			addPlayerPoint.MatchType,
		)
	// Query'yi SQL stringine ve parametrelerine çevir
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("sorgu oluşturulurken hata: %w", err)
	}

	// Sorguyu çalıştır
	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("puan kaydı eklenirken hata (Player: %s): %w", addPlayerPoint.PlayerId, err)
	}

	return nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/service/internal/domain/matchrequest"

	"github.com/Masterminds/squirrel"
)

type MatchRequestRepository struct {
	sqlrepository.Repository
}

func NewMatchRequestRepository(db *sql.DB) *MatchRequestRepository {
	return &MatchRequestRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *MatchRequestRepository) Create(ctx context.Context, dto matchrequest.RequestDto) (string, error) {
	executor := r.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.Insert("match_requests").
		Columns("league_id", "requester_id", "requester_type", "target_date", "start_hour").
		Values(dto.LeagueId, dto.PlayerId, dto.Type, dto.Date, dto.StartHour).
		Suffix("RETURNING id"). // Kaydedilen satırın ID'sini geri almak için
		ToSql()

	if err != nil {
		return "", fmt.Errorf("sorgu oluşturulamadı: %w", err)
	}

	var insertedID string
	// Executor üzerinden sorguyu çalıştırıp dönen ID'yi okuyoruz
	err = executor.QueryRowContext(ctx, query, args...).Scan(&insertedID)
	if err != nil {
		return "", fmt.Errorf("match request kaydedilemedi: %w", err)
	}

	return insertedID, nil
}

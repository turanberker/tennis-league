package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sqlrepository "tennis-league/common/lib/repository/sql"

	"github.com/Masterminds/squirrel"
)

type LeagueCoordinatorRepository struct {
	sqlrepository.Repository
}

func NewLeagueCoordinatorRepository(db *sql.DB) *LeagueCoordinatorRepository {
	return &LeagueCoordinatorRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *LeagueCoordinatorRepository) Exists(ctx context.Context, leagueId string, userId string) (bool, error) {

	executor := r.GetExecutor(ctx)
	query := "SELECT 1 FROM league_coordinator WHERE user_id = $1 AND league_id = $2 LIMIT 1"

	var exists int
	// r.db'nin sql.DB veya sqlx.DB olduğunu varsayıyorum
	err := executor.QueryRowContext(ctx, query, userId, leagueId).Scan(&exists)

	if err != nil {
		// Eğer kayıt bulunamadıysa bu bir hata değil, 'false' durumudur
		if err == sql.ErrNoRows {
			return false, nil
		}
		// Gerçek bir veritabanı hatası varsa (bağlantı vb.) hatayı dönüyoruz
		return false, fmt.Errorf("koordinatör kontrolü yapılamadı: %w", err)
	}

	return true, nil

}

func (r *LeagueCoordinatorRepository) Add(ctx context.Context, leagueId string, userId string) (*bool, error) {
	executor := r.GetExecutor(ctx)

	// SQL sorgusu (PostgreSQL için en performanslı hali)
	query := `INSERT INTO league_coordinator (league_id, user_id) 
              VALUES ($1, $2)
              ON CONFLICT (league_id, user_id) DO NOTHING`

	var result sql.Result
	var err error

	result, err = executor.ExecContext(ctx, query, leagueId, userId)

	if err != nil {
		return nil, fmt.Errorf("koordinatör eklenirken SQL hatası: %w", err)
	}

	// 3. RowsAffected kontrolü:
	// 1 dönerse yeni kayıt eklendi, 0 dönerse ON CONFLICT devreye girdi (kayıt zaten var)
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	isAdded := rows > 0
	return &isAdded, nil
}

func (r *LeagueCoordinatorRepository) GetLeagueIdsByUserId(ctx context.Context, userId string) (*[]string, error) {
	executor := r.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("lc.league_id",
		"lc.league_id").
		From("league_coordinator lc").
		Where(squirrel.Eq{"lc.user_id": userId})

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	// Sorguyu çalıştırıyoruz
	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("sorgu çalıştırılamadı: %v", err)
	}
	defer rows.Close()

	var leagueIds []string
	for rows.Next() {
		var leagueId string
		if err := rows.Scan(&leagueId); err != nil {
			return nil, fmt.Errorf("satır okunamadı: %v", err)
		}
		leagueIds = append(leagueIds, leagueId)
	}

	// Iterasyon sırasında bir hata oluşup oluşmadığını kontrol ediyoruz
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("satırlar işlenirken hata oluştu: %v", err)
	}

	return &leagueIds, nil
}

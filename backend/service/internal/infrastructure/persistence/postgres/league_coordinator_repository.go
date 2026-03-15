package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type LeagueCoordinatorRepository struct {
	db *sql.DB
}

func NewLeagueCoordinatorRepository(db *sql.DB) *LeagueCoordinatorRepository {
	return &LeagueCoordinatorRepository{db: db}
}

func (r *LeagueCoordinatorRepository) Exists(ctx context.Context, leagueId string, userId string) (bool, error) {
	query := "SELECT 1 FROM league_coordinator WHERE user_id = $1 AND league_id = $2 LIMIT 1"

	var exists int
	// r.db'nin sql.DB veya sqlx.DB olduğunu varsayıyorum
	err := r.db.QueryRowContext(ctx, query, userId, leagueId).Scan(&exists)

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
	// 1. Context'ten transaction'ı çekiyoruz
	tx, ok := database.GetTxFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Aktif Transaction yok")

	}

	// SQL sorgusu (PostgreSQL için en performanslı hali)
	query := `INSERT INTO league_coordinator (league_id, user_id) 
              VALUES ($1, $2)
              ON CONFLICT (league_id, user_id) DO NOTHING`

	var result sql.Result
	var err error

	result, err = tx.ExecContext(ctx, query, leagueId, userId)

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

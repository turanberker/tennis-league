package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) GetById(ctx context.Context, id int64) (*player.Player, error) {
	player := &player.Player{}
	query := `SELECT id,  name, surname, sex,user_id FROM players WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&player.ID, &player.Name, &player.Surname, &player.Sex, &player.UserId)
	if err != nil {
		log.Println("Playerı maplerken hata oluştu:", err)
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepository) GetByUuid(ctx context.Context, uuid string) (*player.Player, error) {
	player := &player.Player{}
	query := `SELECT id,  name, surname, sex, user_id FROM players WHERE uuid=$1`
	err := r.db.QueryRowContext(ctx, query, uuid).
		Scan(&player.ID, &player.Name, &player.Surname, &player.Sex, &player.UserId)
	if err != nil {
		log.Println("Playerı maplerken hata oluştu:", err)
		return nil, err
	}
	return player, nil

}
func (r *PlayerRepository) List(ctx context.Context, queryParams player.ListQueryParameters) ([]*player.Player, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("id", "name", "surname", "sex", "user_id").From("players")

	if queryParams.Name != nil {
		sqlBuilder = sqlBuilder.Where("name ILIKE ?", "%"+*queryParams.Name+"%")
	}

	if queryParams.Sex != nil {
		sqlBuilder = sqlBuilder.Where(squirrel.Eq{"sex": *queryParams.Sex})
	}
	if queryParams.HasUser != nil {
		if *queryParams.HasUser {
			sqlBuilder = sqlBuilder.Where("user_id IS NOT NULL")
		} else {
			sqlBuilder = sqlBuilder.Where("user_id IS NULL")
		}
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	type row struct {
		ID      string  `db:"id"`
		Name    string  `db:"name"`
		Surname string  `db:"surname"`
		Sex     string  `db:"sex"`
		UserID  *string `db:"user_id"` // Null gelebileceği için pointer kullanmak güvenlidir
	}
	var rowsData []row
	err = sqlscan.Select(ctx, r.db, &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası: %v", err)
	}

	// Gelen ham veriyi kendi Player modelinize dönüştürme (Mapping)
	var players []*player.Player
	for _, d := range rowsData {
		players = append(players, &player.Player{
			ID:      d.ID,
			Name:    d.Name,
			Surname: d.Surname,
			Sex:     player.Sex(d.Sex),
			UserId:  d.UserID,
		})
	}

	return players, nil
}

func (r *PlayerRepository) Save(ctx context.Context, persistPlayer *player.PersistPlayer) (*string, error) {
	query := `INSERT INTO players (id, name, surname,sex, user_id) VALUES (gen_random_uuid(), $1, $2, $3,$4) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, persistPlayer.Name, persistPlayer.Surname, persistPlayer.Sex, persistPlayer.UserId).
		Scan(&id)
	if err != nil {
		log.Println("Player insert ederken hata oluştu:", err)
		return nil, err
	}
	return &id, nil
}

func (r *PlayerRepository) AssignToUser(ctx context.Context, playerId string, userId string) error {

	tx, ok := database.GetTxFromContext(ctx)
	if !ok {
		panic("Aktif Transaction yok")
	}
	query := `UPDATE players SET user_id = $1 WHERE id = $2`

	result, err := tx.ExecContext(ctx, query, userId, playerId)
	if err != nil {
		return fmt.Errorf("Oyuncuya kullanıcı ataması başarısız:")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return fmt.Errorf("Güncellenecek oyuncu bulunamadı: %w", err)
	}

	return nil

}

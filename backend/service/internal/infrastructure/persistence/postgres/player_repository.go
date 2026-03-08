package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

	query := `SELECT id, name, surname, sex, user_id FROM players WHERE 1=1 and `
	args := []interface{}{}
	if queryParams.Name != nil {
		query += query + " and name ILIKE '%' || $1 || '%'"
		args = append(args, queryParams.Name)
	}

	if queryParams.Sex != nil {
		query = query + " and sex $2"
		args = append(args, queryParams.Sex)
	}
	if queryParams.HasUser != nil {
		if *queryParams.HasUser {
			query += " AND user_id IS NOT NULL"
		} else {
			query += " AND user_id IS NULL"
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("Player listesi çekerken hata oluştu:", err)
		return nil, err
	}
	defer rows.Close()
	var players []*player.Player
	for rows.Next() {
		player := &player.Player{}
		err := rows.Scan(&player.ID, &player.Name, &player.Surname, &player.Sex, &player.UserId)

		if err != nil {
			log.Println("Playerları maplerken hata oluştu:", err)
			return nil, err
		}
		players = append(players, player)
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

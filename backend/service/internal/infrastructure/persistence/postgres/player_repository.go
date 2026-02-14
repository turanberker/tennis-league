package postgres

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/domain/player"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) GetById(ctx context.Context, id int64) (*player.Player, error) {
	player := &player.Player{}
	query := `SELECT id, uuid, name, surname, user_id FROM players WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&player.ID, &player.Uuid, &player.Name, &player.Surname, &player.UserId)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepository) GetByUuid(ctx context.Context, uuid string) (*player.Player, error) {
	player := &player.Player{}
	query := `SELECT id, uuid, name, surname, user_id FROM players WHERE uuid=$1`
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(&player.ID, &player.Uuid, &player.Name, &player.Surname, &player.UserId)
	if err != nil {
		return nil, err
	}
	return player, nil

}
func (r *PlayerRepository) List(ctx context.Context, name string) ([]*player.Player, error) {
	query := `SELECT id, uuid, name, surname, user_id FROM players WHERE name ILIKE '%' || $1 || '%'`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []*player.Player
	for rows.Next() {
		player := &player.Player{}
		err := rows.Scan(&player.ID, &player.Uuid, &player.Name, &player.Surname, &player.UserId)

		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (r *PlayerRepository) Save(ctx context.Context, persistPlayer *player.PersistPlayer) (int64, error) {
	query := `INSERT INTO players (uuid, name, surname, user_id) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, persistPlayer.Name, persistPlayer.Surname, persistPlayer.UserId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

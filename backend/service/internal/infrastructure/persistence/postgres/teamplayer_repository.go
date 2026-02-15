package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/domain/teamplayer"
)

type TeamPlayerRepository struct {
	db *sql.DB
}

func NewTeamPlayerRepository(db *sql.DB) *TeamPlayerRepository {
	return &TeamPlayerRepository{db: db}
}

func (r *TeamPlayerRepository) GetByPlayersByTeamId(ctx context.Context, teamId string) ([]*player.Player, error) {
	query := `SELECT  p.id, p.name,p.surname,p.sex,p.user_id FROM team_players tp inner join players p on p.id=tp.player_id WHERE team_id=$1`
	rows, err := r.db.QueryContext(ctx, query, teamId)

	if err != nil {
		log.Println("Takım Oyuncuları çekilirken hata oluştu:", err)
		return nil, err
	}
	defer rows.Close()
	var teamPlayers []*player.Player
	for rows.Next() {
		teamPlayer := &player.Player{}
		if err := rows.Scan(
			&teamPlayer.ID,
			&teamPlayer.Name,
			&teamPlayer.Surname,
			&teamPlayer.Sex,
			&teamPlayer.UserId,
		); err != nil {
			return nil, err
		}

		teamPlayers = append(teamPlayers, teamPlayer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teamPlayers, nil
}

func (r *TeamPlayerRepository) Save(ctx context.Context, tx *sql.Tx, teamPlayer *teamplayer.PersistTeamPlayer) error {
	query := `INSERT INTO team_players ( team_id, player_id) VALUES ($1, $2) `

	_, err := tx.ExecContext(ctx, query, teamPlayer.TeamID, teamPlayer.PlayerID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {

			// P0001 = RAISE EXCEPTION
			if pqErr.Code == "P0001" {
				return errors.New(pqErr.Message)
			}
		} else {
			log.Println("Takım Oyuncusu insert ederken hata oluştu:", err)
			return err
		}

	}
	return nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/service/internal/domain/player"
	"tennis-league/service/internal/domain/teamplayer"

	"github.com/lib/pq"
)

type TeamPlayerRepository struct {
	sqlrepository.Repository
}

func NewTeamPlayerRepository(db *sql.DB) *TeamPlayerRepository {
	return &TeamPlayerRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *TeamPlayerRepository) GetByPlayersByTeamId(ctx context.Context, teamId string) ([]*player.Player, error) {
	exec := r.GetExecutor(ctx)
	query := `SELECT  p.id, p.name, p.surname, p.sex, p.user_id,p.single_point ,p.double_point FROM team_player tp inner join player p on p.id=tp.player_id WHERE team_id=$1`
	rows, err := exec.QueryContext(ctx, query, teamId)

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
			&teamPlayer.SinglePoints,
			&teamPlayer.DoublePoints,
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

func (r *TeamPlayerRepository) Save(ctx context.Context, teamPlayer *teamplayer.PersistTeamPlayer) error {
	exec := r.GetExecutor(ctx)
	query := `INSERT INTO team_player ( team_id, player_id) VALUES ($1, $2) `

	_, err := exec.ExecContext(ctx, query, teamPlayer.TeamID, teamPlayer.PlayerID)
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

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
)

type ScoreBoardRepository struct {
	db *sql.DB
}

func NewScoreBoardRepository(db *sql.DB) *ScoreBoardRepository {
	return &ScoreBoardRepository{db: db}
}

func (f *ScoreBoardRepository) SaveFixture(ctx context.Context, tx *sql.Tx, leagueId string, teams []string) error {

	if len(teams) == 0 {
		return errors.New("Takım Listesi boş olamaz")
	}

	// INSERT INTO fixtures (league_id, home_team_id, away_team_id) VALUES ($1,$2,$3), ...
	query := "INSERT INTO score_board (league_id, team_id) VALUES "
	args := []interface{}{}

	for i, m := range teams {
		// Her match için 2 parametre
		query += fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2)
		if i != len(teams)-1 {
			query += ", "
		}
		args = append(args, leagueId, m)
	}

	// Bulk insert
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("Scoreboard initlenirken hata oluştu:", err)
		return err
	}

	return nil
}

func (f *ScoreBoardRepository) GetScoreBoard(ctx context.Context, leagueId string) ([]*scoreboard.ScoreBoard, error) {
	query := `Select f.team_id,t.name ,f.played,f.won, f.lost,
		f.won_sets, f.lost_sets,f.won_games,f.lost_games,score
		from score_board f 
		inner join teams  t on t.id =f.team_id 
		where f.league_id=$1 order by score desc`

	rows, err := f.db.QueryContext(ctx, query, leagueId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fixtures []*scoreboard.ScoreBoard

	for rows.Next() {
		team := &scoreboard.ScoreBoard{}
		err := rows.Scan(&team.Team.Id, &team.Team.Name, &team.Played, &team.Won, &team.Lost,
			&team.WonSets, &team.LostSets, &team.WonGames, &team.LostGames, &team.Score)

		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, team)
	}
	return fixtures, nil

}

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

func (f *ScoreBoardRepository) UpdateScore(ctx context.Context, tx *sql.Tx, update scoreboard.IncreaseTeamScore) error {
	query := `
		UPDATE score_board
		SET 
			played = played+1,
			won_sets = won_sets + $1,
			lost_sets = lost_sets + $2,
			won_games = won_games + $3,
			lost_games = lost_games + $4,
			score = score + $5,
			won = won + CASE WHEN $6 THEN 1 ELSE 0 END,
			lost = lost + CASE WHEN $6 THEN 0 ELSE 1 END
		WHERE league_id = $7
		  AND team_id = $8
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		update.WonSets,
		update.LostSets,
		update.WonGames,
		update.LostGames,
		update.IncreaseScore,
		update.Won,
		update.LeagueId,
		update.TeamId,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("scoreboard not found for team %s", update.TeamId)
	}

	return nil
}

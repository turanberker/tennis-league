package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	matchSet "github.com/turanberker/tennis-league-service/internal/domain/matchset"
)

type MatchSetRepository struct {
	db *sql.DB
}

func NewMatchSetRepository(db *sql.DB) *MatchSetRepository {
	return &MatchSetRepository{db: db}
}

func (r *MatchSetRepository) SaveSetScore(ctx context.Context, tx *sql.Tx, setScore *matchSet.UpdateSetScore) error {
	query := `INSERT INTO match_sets ( match_id, set_number,team_1_games, team_2_games) VALUES ($1, $2,$3,$4)`

	_, err := tx.ExecContext(ctx, query, setScore.MatchId, setScore.Set, setScore.Team1Score, setScore.Team2Score)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {

			// P0001 = RAISE EXCEPTION
			if pqErr.Code == "P0001" {
				return errors.New(pqErr.Message)
			}
		} else {
			log.Printf("%d skoru insert edilirken hata oluştu: %v", setScore.Set, err)
			return err
		}

	}
	return nil
}

func (r *MatchSetRepository) SaveSuperTieScore(ctx context.Context, tx *sql.Tx, setScore *matchSet.UpdateSuperTieScore) error {
	query := `INSERT INTO match_sets ( match_id, set_number,team_1_tie_break_score, team_2_tie_break_score) VALUES ($1, 3,$2,$3)`
	_, err := tx.ExecContext(ctx, query, setScore.MatchId, setScore.Team1Score, setScore.Team2Score)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {

			// P0001 = RAISE EXCEPTION
			if pqErr.Code == "P0001" {
				return errors.New(pqErr.Message)
			}
		} else {
			log.Println("Super Tie skoru insert edilirken hata oluştu:", err)
			return err
		}

	}
	return nil
}

func (r *MatchSetRepository) DeleteSetScores(ctx context.Context, tx *sql.Tx, matchId string) error {
	query := "delete from match_sets where match_id =$1"

	_, err := tx.ExecContext(ctx, query, matchId)
	if err != nil {
		log.Printf("%s id li mac skorları silinemedi", matchId)
		return err
	}
	return nil
}

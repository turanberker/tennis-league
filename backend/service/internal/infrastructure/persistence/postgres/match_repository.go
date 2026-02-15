package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/turanberker/tennis-league-service/internal/domain/match"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) SaveLeagueMatches(ctx context.Context, tx *sql.Tx, matches []match.PersistLeagueMatch) error {
	if len(matches) == 0 {
		return nil // insert edilecek maç yok
	}

	// INSERT INTO fixtures (league_id, home_team_id, away_team_id) VALUES ($1,$2,$3), ...
	query := "INSERT INTO tennisleague.matches (league_id, team_1_id, team_2_id) VALUES "
	args := []interface{}{}

	for i, m := range matches {
		// Her match için 3 parametre
		query += fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3)
		if i != len(matches)-1 {
			query += ", "
		}
		args = append(args, m.LeagueId, m.Team1Id, m.Team2Id)
	}

	// Bulk insert
	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("League fixture'ları insert edilirken hata oluştu:", err)
		return err
	}

	return nil
}

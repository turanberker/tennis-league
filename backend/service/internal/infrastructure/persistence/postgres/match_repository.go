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

func (r *MatchRepository) SaveLeagueMatches(ctx context.Context, tx *sql.Tx, matches []*match.PersistLeagueMatch) error {
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

func (r *MatchRepository) GetFixtureByLeagueId(ctx context.Context, leagueId string) ([]*match.LeagueFixtureMatch, error) {
	query := `
		SELECT m.id, m.team_1_id, t1.name, m.team_2_id, t2.name, m.status, m.match_date
		FROM tennisleague.matches m
		JOIN tennisleague.teams t1 ON m.team_1_id = t1.id
		JOIN tennisleague.teams t2 ON m.team_2_id = t2.id
		WHERE m.league_id = $1 order by m.match_date asc
	`
	rows, err := r.db.QueryContext(ctx, query, leagueId)
	if err != nil {
		log.Println("Maç listesi çekerken hata oluştu:", err)
		return nil, err
	}

	defer rows.Close()
	var matches []*match.LeagueFixtureMatch
	for rows.Next() {
		match := &match.LeagueFixtureMatch{}
		err := rows.Scan(&match.Id, &match.Team1.Id, &match.Team1.Name, &match.Team2.Id, &match.Team2.Name, &match.Status, &match.MatchDate)

		if err != nil {
			log.Println("Maçlar maplerken hata oluştu:", err)
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

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
		SELECT m.id, m.team_1_id, t1.name,m.team_1_score, m.winner_id =m.team_1_id, 
		m.team_2_id, t2.name,m.team_2_score ,m.winner_id =m.team_2_id, m.status, m.match_date
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
		err := rows.Scan(&match.Id,
			&match.Team1.Id, &match.Team1.Name, &match.Team1.Score, &match.Team1.Winner,
			&match.Team2.Id, &match.Team2.Name, &match.Team2.Score, &match.Team2.Winner,
			&match.Status, &match.MatchDate)

		if err != nil {
			log.Println("Maçlar maplerken hata oluştu:", err)
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *MatchRepository) UpdateMatchDate(ctx context.Context, tx *sql.Tx, data match.UpdateMatchDate) error {
	query := "Update matches set match_date=$1 where id=$2"

	_, err := tx.ExecContext(ctx, query, data.MatchDate, data.Id)
	if err != nil {
		log.Println("Maç tarihi güncellenirken hata oluştu:", err)
		return err
	}
	return nil
}

func (r *MatchRepository) GetMatchTeamIds(ctx context.Context, matchId string) *match.MatchTeamIds {

	var response match.MatchTeamIds
	query := "select team_1_id ,team_2_id ,status  from matches m where id=$1"

	err := r.db.QueryRowContext(ctx, query, matchId).
		Scan(&response.Team1Id, &response.Team2Id, &response.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			// kayıt yok
			return nil
		}
		return nil
	}
	return &response
}

func (r *MatchRepository) UpdateMatchScore(ctx context.Context, tx *sql.Tx, macScore *match.UpdateMatchScore) error {
	query := "Update matches set team_1_score=$1, team_2_score=$2, winner_id=$3, status=$4 where id=$5"

	_, err := tx.ExecContext(ctx, query, macScore.Team1Score, macScore.Team2Score, macScore.WinnerTeamId, match.StatusCompleted, macScore.Id)
	if err != nil {
		log.Printf("Maç Skoru güncellenirken hata oluştu:%+v", err)
		return err
	}
	return nil
}

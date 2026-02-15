package postgres

import (
	context "context"
	"database/sql"
	"log"

	"github.com/turanberker/tennis-league-service/internal/domain/team"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) GetById(ctx context.Context, id string) (*team.Team, error) {
	team := &team.Team{}
	query := `SELECT id, league_id, name FROM teams WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&team.ID, &team.LeagueID, &team.Name)
	return team, err
}

func (r *TeamRepository) GetByLeagueId(ctx context.Context, leagueId string) ([]*team.Team, error) {
	query := `SELECT id, league_id, name FROM teams WHERE league_id=$1`
	rows, err := r.db.QueryContext(ctx, query, leagueId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []*team.Team
	for rows.Next() {
		team := &team.Team{}
		err := rows.Scan(&team.ID, &team.LeagueID, &team.Name)

		if err != nil {
			log.Println("Takımlar maplerken hata oluştu:", err)
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) Save(ctx context.Context, tx *sql.Tx, persistTeam *team.PersistTeam) (*string, error) {
	query := `INSERT INTO teams (id, league_id, name) VALUES (gen_random_uuid(), $1, $2) RETURNING id`
	var id string
	err := tx.QueryRowContext(ctx, query, persistTeam.LeagueID, persistTeam.Name).Scan(&id)
	if err != nil {
		log.Println("Takım insert ederken hata oluştu:", err)
		return nil, err
	}
	return &id, nil
}

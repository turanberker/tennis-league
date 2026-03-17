package postgres

import (
	context "context"
	"database/sql"
	"log"

	"github.com/turanberker/tennis-league-service/internal/domain/team"
)

type TeamRepository struct {
	BaseRepository
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{BaseRepository{db: db}}
}

func (r *TeamRepository) GetById(ctx context.Context, id string) (*team.Team, error) {

	exec := r.GetExecutor(ctx)
	team := &team.Team{}
	query := `SELECT id, league_id, name FROM team WHERE id=$1`
	err := exec.QueryRowContext(ctx, query, id).Scan(&team.ID, &team.LeagueID, &team.Name)
	return team, err
}

func (r *TeamRepository) GetByLeagueId(ctx context.Context, leagueId string) ([]*team.Team, error) {
	exec := r.GetExecutor(ctx)
	query := `SELECT id, league_id, name FROM team WHERE league_id=$1`
	rows, err := exec.QueryContext(ctx, query, leagueId)
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

func (r *TeamRepository) Save(ctx context.Context, persistTeam *team.PersistTeam) (*string, error) {
	exec := r.GetExecutor(ctx)
	query := `INSERT INTO team (league_id, name) VALUES ($1, $2) RETURNING id`
	var id string
	err := exec.QueryRowContext(ctx, query, persistTeam.LeagueID, persistTeam.Name).Scan(&id)
	if err != nil {
		log.Println("Takım insert ederken hata oluştu:", err)
		return nil, err
	}
	return &id, nil
}

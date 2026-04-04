package postgres

import (
	context "context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
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

func (r *TeamRepository) GetByLeagueId(ctx context.Context, leagueId string) ([]*team.LeagueTeam, error) {
	exec := r.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("t.id", "t.name", "sum(p.double_point) power").
		From("team t").InnerJoin("team_player tp on tp.team_id = t.id").
		InnerJoin("player p on p.id =tp.player_id ").
		Where(squirrel.Eq{"t.league_id": leagueId}).
		GroupBy("t.id", "t.name").OrderBy("power DESC")

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	type row struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		POWER int32  `db:"power"`
	}
	var rowsData []row
	err = sqlscan.Select(ctx, exec, &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası: %v", err)
	}

	var teams []*team.LeagueTeam
	for _, d := range rowsData {
		teams = append(teams, &team.LeagueTeam{
			ID:    d.ID,
			Name:  d.Name,
			Power: d.POWER,
		})
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

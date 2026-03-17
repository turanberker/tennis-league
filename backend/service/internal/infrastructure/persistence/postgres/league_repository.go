package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/lib/pq"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
)

type LeagueRepository struct {
	BaseRepository
}

func NewLeagueRepository(db *sql.DB) *LeagueRepository {
	return &LeagueRepository{BaseRepository: BaseRepository{db: db}}
}

func (r *LeagueRepository) GetById(ctx context.Context, id string) (*league.League, error) {
	exec := r.GetExecutor(ctx)
	league := &league.League{}
	query := `SELECT id, name,fixture_created_date FROM league WHERE id=$1`
	err := exec.QueryRowContext(ctx, query, id).Scan(&league.ID, &league.Name, &league.FixtureCreatedDate)
	if err != nil {
		return nil, err
	}
	return league, err
}

func (r *LeagueRepository) GetAll(ctx context.Context, name *string) ([]*league.League, error) {
	exec := r.GetExecutor(ctx)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sqlBuilder := psql.Select("l.id",
		"l.name",
		"l.fixture_created_date",
		"COALESCE(string_agg(concat(u.name, ' ', u.surname), ','),'') AS coordinators",
		"COALESCE(string_agg(u.id,','),'') as coordinator_user_ids ").
		From("league l").
		LeftJoin("league_coordinator lc ON lc.league_id = l.id").
		LeftJoin("\"user\" u ON u.id = lc.user_id").
		GroupBy("l.id", "l.name", "l.fixture_created_date")
	if name != nil && len(*name) > 0 {
		sqlBuilder = sqlBuilder.Where("l.name ILIKE ?", "%"+*name+"%")
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	type row struct {
		ID                 string     `db:"id"`
		Name               string     `db:"name"`
		FixtureCreatedDate *time.Time `db:"fixture_created_date"`
		Coordinators       string     `db:"coordinators"`
		CoordinatorIds     string     `db:"coordinator_user_ids"`
	}
	var rowsData []row
	err = sqlscan.Select(ctx, exec, &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası: %v", err)
	}

	var leagues []*league.League
	for _, d := range rowsData {

		leagues = append(leagues, &league.League{
			ID:                 d.ID,
			Name:               d.Name,
			FixtureCreatedDate: d.FixtureCreatedDate,
			Cootrinators:       strings.Split(d.Coordinators, ","),
			CoordinatorUserId:  strings.Split(d.CoordinatorIds, ","),
		})
	}

	return leagues, nil
}

func (r *LeagueRepository) Save(ctx context.Context, persistLeague *league.PersistLeague) (*string, error) {
	exec := r.GetExecutor(ctx)
	query := `INSERT INTO league (id,name) VALUES (gen_random_uuid(),$1) RETURNING id`

	var id string
	err := exec.QueryRowContext(ctx, query, persistLeague.Name).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "league_name_key" {
			return nil, league.LEAGE_WITH_NAME_EXISTS
		}

		return nil, err
	}
	return &id, nil
}

func (r *LeagueRepository) SetFitxtureCreatedDate(ctx context.Context, leagueId string) error {
	exec := r.GetExecutor(ctx)
	query := `UPDATE league SET fixture_created_date = NOW() WHERE id = $1`
	_, err := exec.ExecContext(ctx, query, leagueId)
	return err
}

func (r *LeagueRepository) IsFixtureCreated(ctx context.Context, leagueId string) (bool, error) {
	exec := r.GetExecutor(ctx)
	query := `SELECT EXISTS (
		SELECT 1 
		FROM league 
		WHERE id = $1 
		  AND fixture_created_date IS NOT NULL
	)`
	var exists bool
	err := exec.QueryRowContext(ctx, query, leagueId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

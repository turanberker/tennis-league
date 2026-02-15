package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
)

type LeagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *sql.DB) *LeagueRepository {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) GetById(ctx context.Context, id int64) (*league.League, error) {
	league := &league.League{}
	query := `SELECT id, name,fixture_created_date FROM leagues WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&league.ID, &league.Name, &league.FixtureCreatedDate)
	if err != nil {
		return nil, err
	}
	return league, err
}

func (r *LeagueRepository) GetAll(ctx context.Context, name string) ([]*league.League, error) {
	query := `SELECT id, name, fixture_created_date FROM leagues WHERE name ILIKE '%' || $1 || '%'`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var leagues []*league.League

	for rows.Next() {
		league := &league.League{}
		err := rows.Scan(&league.ID, &league.Name, &league.FixtureCreatedDate)

		if err != nil {
			return nil, err
		}
		leagues = append(leagues, league)

	}
	return leagues, nil
}

func (r *LeagueRepository) Save(ctx context.Context, persistLeague *league.PersistLeague) (*string, error) {
	query := `INSERT INTO leagues (id,name) VALUES (gen_random_uuid(),$1) RETURNING id`

	var id string
	err := r.db.QueryRowContext(ctx, query, persistLeague.Name).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "leagues_name_key" {
			return nil, errors.New("league name already exists")
		}

		return nil, err
	}
	return &id, nil
}

func (r *LeagueRepository) SetFitxtureCreatedDate(ctx context.Context, tx *sql.Tx, leagueId string) error {
	query := `UPDATE leagues SET fixture_created_date = NOW() WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, leagueId)
	return err
}

func (r *LeagueRepository) IsFixtureCreated(ctx context.Context, leagueId string) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 
		FROM leagues 
		WHERE id = $1 
		  AND fixture_created_date IS NOT NULL
	)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, leagueId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

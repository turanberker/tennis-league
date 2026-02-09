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
	query := `SELECT id, name FROM leagues WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&league.ID, &league.Name)
	if err != nil {
		return nil, err
	}
	return league, err
}

func (r *LeagueRepository) GetAll(ctx context.Context, name string) ([]*league.League, error) {
	query := `SELECT id, name FROM leagues WHERE name ILIKE '%' || $1 || '%'`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var leagues []*league.League

	for rows.Next() {
		league := &league.League{}
		err := rows.Scan(&league.ID, &league.Name)

		if err != nil {
			return nil, err
		}
		leagues = append(leagues, league)

	}
	return leagues, nil
}

func (r *LeagueRepository) Save(ctx context.Context, persistLeague *league.PersistLeague) (int64, error) {
	query := `INSERT INTO leagues (name) VALUES ($1) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, persistLeague.Name).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "leagues_name_key" {
			return 0, errors.New("league name already exists")
		}

		return 0, err
	}
	return id, nil
}

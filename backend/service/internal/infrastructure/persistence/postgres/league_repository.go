package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("l.id",
		"l.name",
		"l.format",
		"l.category",
		"l.process_type",
		"l.status",
		"l.total_attendance",
		"l.start_date",
		"l.end_date").
		From("league l").
		Where(squirrel.Eq{"l.id": id})

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}
	league := &league.League{}

	if err := sqlscan.Get(ctx, exec, league, query, args...); err != nil {
		return nil, fmt.Errorf("lig getirilemedi (id: %s): %w", id, err)
	}

	return league, nil
}

func (r *LeagueRepository) GetAll(ctx context.Context, stauts *league.LEAGUE_STATUS) ([]*league.LeagueListSelect, error) {
	exec := r.GetExecutor(ctx)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sqlBuilder := psql.Select("l.id",
		"l.name",
		"l.format",
		"l.category",
		"l.process_type",
		"l.status",
		"l.total_attendance",
		"COALESCE(string_agg(lc.user_id,','),'') as coordinator_user_ids").
		From("league l").
		LeftJoin("league_coordinator lc ON lc.league_id = l.id").
		GroupBy("l.id", "l.name", "l.format", "l.process_type", "l.status", "l.total_attendance")
	if stauts != nil {
		sqlBuilder = sqlBuilder.Where(squirrel.Eq{"l.status": *stauts})
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	type row struct {
		ID               string                     `db:"id"`
		Name             string                     `db:"name"`
		Format           league.LEAGUE_FORMAT       `db:"format"`
		Category         league.LEAGUE_CATEGORY     `db:"category"`
		Status           league.LEAGUE_STATUS       `db:"status"`
		Process_Type     league.LEAGUE_PROCESS_TYPE `db:"process_type"`
		Total_Attancande int32                      `db:"total_attendance"`
		CoordinatorIds   string                     `db:"coordinator_user_ids"`
	}
	var rowsData []row
	err = sqlscan.Select(ctx, exec, &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası: %v", err)
	}

	var leagues []*league.LeagueListSelect
	for _, d := range rowsData {

		leagues = append(leagues, &league.LeagueListSelect{
			ID:                d.ID,
			Name:              d.Name,
			Category:          d.Category,
			Format:            d.Format,
			Type:              d.Process_Type,
			Status:            d.Status,
			TotalAttentance:   d.Total_Attancande,
			CoordinatorUserId: strings.Split(d.CoordinatorIds, ","),
		})
	}

	return leagues, nil
}

func (r *LeagueRepository) Save(ctx context.Context, persistLeague *league.PersistLeague) (*string, error) {
	exec := r.GetExecutor(ctx)
	query := `INSERT INTO league (name,format, category,process_type) VALUES ($1,$2,$3,$4) RETURNING id`

	var id string
	err := exec.QueryRowContext(ctx, query, persistLeague.Name, persistLeague.Format, persistLeague.Categoty, persistLeague.ProcessType).
		Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Constraint == "league_name_key" {
			return nil, league.LEAGE_WITH_NAME_EXISTS
		}

		return nil, err
	}
	return &id, nil
}

func (r *LeagueRepository) StartLeague(ctx context.Context, leagueId string) error {
	exec := r.GetExecutor(ctx)
	query := `UPDATE league SET start_date = NOW(), status=$1 WHERE id = $2`
	_, err := exec.ExecContext(ctx, query, league.LeagueStatus_ACTIVE, leagueId)
	return err
}

func (r *LeagueRepository) IsFixtureCreated(ctx context.Context, leagueId string) (bool, error) {
	exec := r.GetExecutor(ctx)
	query := `SELECT EXISTS (
		SELECT 1 
		FROM league 
		WHERE id = $1 
		  AND status != $2
	)`
	var exists bool
	err := exec.QueryRowContext(ctx, query, leagueId, league.LeagueStatus_DRAFT).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *LeagueRepository) IncreaseAttandanceCount(ctx context.Context, leagueId string) (*int32, error) {
	exec := r.GetExecutor(ctx)
	query := `UPDATE league SET total_attendance =total_attendance+1  WHERE id = $1 RETURNING total_attendance`
	var updatedAttendance int32
	err := exec.QueryRowContext(ctx, query, leagueId).Scan(&updatedAttendance)

	if err != nil {
		return nil, err
	}

	return &updatedAttendance, nil
}

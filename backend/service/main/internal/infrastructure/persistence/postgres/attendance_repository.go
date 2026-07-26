package postgres

import (
	context "context"
	"database/sql"

	"fmt"
	"log"
	"net/http"

	customerror "tennis-league/common/lib/error"
	sqlrepository "tennis-league/common/lib/repository/sql"
	errorcodes "tennis-league/service/internal/domain/error_codes"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/team"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type AttendanceRepository struct {
	sqlrepository.Repository
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *AttendanceRepository) GetByLeagueId(ctx context.Context, leagueId string) ([]*team.LeagueTeam, error) {
	exec := r.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("t.id", "t.name", "sum(p.double_point) power").
		From("attendance t").InnerJoin("attendance_player tp on tp.attendance_id = t.id").
		InnerJoin("player p on p.id =tp.player_id ").
		Where(squirrel.Eq{"t.league_id": leagueId}).
		GroupBy("t.id", "t.name").OrderBy("power DESC")

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sorgu oluşturulamadı")
	}

	type row struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		POWER int32  `db:"power"`
	}
	var rowsData []row
	err = sqlscan.Select(ctx, exec, &rowsData, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "veritabanı hatası")
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

func (r *AttendanceRepository) Save(ctx context.Context, persistTeam *team.PersistTeam) (*string, error) {
	exec := r.GetExecutor(ctx)

	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("league_attendance").
		Columns("league_id", "name").
		Values(persistTeam.LeagueID, persistTeam.Name).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return nil, customerror.NewInternalError(errors.Wrap(err, "sorgu oluşturulamadı"))
	}

	var id string
	err = exec.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		// unique constraint hatası gibi özel durumları burada da kontrol edebilirsiniz.
		log.Println("Takım insert ederken hata oluştu:", err)
		return nil, customerror.NewInternalError(err)
	}
	return &id, nil
}

func (f *AttendanceRepository) AddPlayerToLeague(ctx context.Context, leagueAttendance league.NewLeaguePlayerAttendance) error {
	executor := f.GetExecutor(ctx)
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("attendance").
		Columns("league_id", "name").
		Values(leagueAttendance.LeagueId, leagueAttendance.PlayerName).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return customerror.NewInternalError(err)
	}

	var id string
	err = executor.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Constraint == "uidx_league_single_player" {
				return customerror.NewBusinessError(http.StatusConflict, errorcodes.ErrorPlayerAlreadyAddedToLeague, "Bu oyuncu bu lige zaten eklenmiş.")
			}
		}

		// Veritabanı hatasını projenizin özel hata yapısıyla sarmalıyoruz
		return customerror.NewInternalError(err)
	}

	query, args, err = squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("attendance_player").
		Columns("attendance_id", "player_id").
		Values(id, leagueAttendance.PlayerId).
		ToSql()
	if err != nil {
		return customerror.NewInternalError(err)
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Constraint == "uidx_league_single_player" {
				return customerror.NewBusinessError(http.StatusConflict, errorcodes.ErrorPlayerAlreadyAddedToLeague, "Bu oyuncu bu lige zaten eklenmiş.")
			}
		}

		// Veritabanı hatasını projenizin özel hata yapısıyla sarmalıyoruz
		return customerror.NewInternalError(err)
	}

	// Her şey yolunda gittiyse nil dönüyoruz
	return nil
}

func (f *AttendanceRepository) IsPlayerAttendedToLeague(ctx context.Context, leagueId string, playerId string) (bool, error) {
	executor := f.GetExecutor(ctx)

	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"1",
		).
		From("attendance_player p").
		InnerJoin("attendance a ON a.id = p.attendance_id").
		Where(squirrel.Eq{
			"a.league_id": leagueId,
			"p.player_id": playerId,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("Sorgu oluşturulamadı: %w", err)
	}

	var dummy int

	err = executor.QueryRowContext(ctx, query, args...).Scan(&dummy)
	if err != nil {
		// Eğer hiç kayıt bulunamadıysa sql.ErrNoRows döner, bu bir hata değil oyuncunun ligde olmadığını gösterir.
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("Sorcu çalıştırılamadı: %w", err)
	}

	// Kayıt başarıyla scan edildiyse oyuncu lige katılmıştır.
	return true, nil
}

func (f *AttendanceRepository) SingleLeagueAttendanceList(ctx context.Context, leagueId string) ([]league.SingleLeagueAttendance, error) {
	executor := f.GetExecutor(ctx)

	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(
			"sb.player_id",
			"p.name",
			"p.surname",
			"p.single_point",
		).
		From("attendance a").
		Join("attendance_player sb ON sb.attendance_id = a.id").
		Join("player p ON p.id = sb.player_id").
		Where(squirrel.Eq{
			"a.league_id": leagueId,
		}).
		ToSql()
	if err != nil {
		return nil, customerror.NewInternalError(err)
	}

	type attendanceRow struct {
		ID           string `db:"player_id"`
		Name         string `db:"name"`
		Surname      string `db:"surname"`
		SinglePoints int    `db:"single_point"`
	}
	var rowsData []attendanceRow
	err = sqlscan.Select(ctx, executor, &rowsData, query, args...)
	if err != nil {
		return nil, customerror.NewInternalError(err)
	}

	// Gelen ham veriyi kendi Player modelinize dönüştürme (Mapping)
	players := make([]league.SingleLeagueAttendance, 0, len(rowsData))
	for _, d := range rowsData {
		players = append(players, league.SingleLeagueAttendance{
			ID:        d.ID,
			Firstname: d.Name,
			Surname:   d.Surname,
			Power:     d.SinglePoints,
		})
	}

	return players, nil
}

func (f *AttendanceRepository) AddTeamToLeague(ctx context.Context, leagueId string, teamId string) error {
	executor := f.GetExecutor(ctx)

	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("score_board").
		Columns("league_id", "team_id").
		Values(leagueId, teamId).
		ToSql()
	if err != nil {
		return customerror.NewInternalError(err)
	}
	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Constraint == "uidx_league_team" {
				return customerror.NewBusinessError(http.StatusConflict, errorcodes.ErrorPlayerAlreadyAddedToLeague, "Bu takım bu lige zaten eklenmiş.")
			}
		}

		// Veritabanı hatasını projenizin özel hata yapısıyla sarmalıyoruz
		return customerror.NewInternalError(err)
	}

	// Her şey yolunda gittiyse nil dönüyoruz
	return nil
}

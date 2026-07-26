package fixturerepository

import (
	"context"
	"database/sql"
	"fmt"

	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/service/internal/delivery/message/consumer/match_score/leaguematch"
	"tennis-league/service/internal/domain/scoreboard"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/pkg/errors"
)

type ScoreBoardRepository struct {
	sqlrepository.Repository
}

func NewScoreBoardRepository(db *sql.DB) *ScoreBoardRepository {
	return &ScoreBoardRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (f *ScoreBoardRepository) FetchScoreBoard(ctx context.Context, leagueId string) ([]*scoreboard.ScoreBoard, error) {
	exec := f.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select(
		"f.attendance_id",
		"t.name",
		"f.played",
		"f.won",
		"f.lost",
		"f.won_sets",
		"f.lost_sets",
		"f.won_games",
		"f.lost_games",
		"f.score",
	).
		From("score_board f").
		InnerJoin("attendance t ON t.id = f.attendance_id").
		Where(squirrel.Eq{"f.league_id": leagueId}).
		OrderBy("f.score DESC")

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "puan durumu sorgusu oluşturulamadı")
	}

	// Düz ve basit bir DTO (Data Transfer Object) tanımlıyoruz.
	// Bu, veritabanı kolonlarıyla doğrudan eşleşir.
	type scoreBoardRow struct {
		AttendanceId string `db:"attendance_id"`
		Name         string `db:"name"`
		Played       int16  `db:"played"`
		Won          int16  `db:"won"`
		Lost         int16  `db:"lost"`
		WonSets      int16  `db:"won_sets"`
		LostSets     int16  `db:"lost_sets"`
		WonGames     int16  `db:"won_games"`
		LostGames    int16  `db:"lost_games"`
		Score        int16  `db:"score"`
	}

	var rowsData []*scoreBoardRow
	if err := sqlscan.Select(ctx, exec, &rowsData, query, args...); err != nil {
		return nil, errors.Wrap(err, "puan durumu getirilirken veritabanı hatası")
	}

	// DTO'dan domain modeline dönüşüm yapıyoruz.
	scoreBoards := make([]*scoreboard.ScoreBoard, 0, len(rowsData))
	for _, row := range rowsData {
		scoreBoards = append(scoreBoards, &scoreboard.ScoreBoard{
			Team: scoreboard.AttendanceReferance{
				Id:   row.AttendanceId,
				Name: row.Name,
			},
			Played:    row.Played,
			Won:       row.Won,
			Lost:      row.Lost,
			WonSets:   row.WonSets,
			LostSets:  row.LostSets,
			WonGames:  row.WonGames,
			LostGames: row.LostGames,
			Score:     row.Score,
		})
	}

	return scoreBoards, nil
}

func (f *ScoreBoardRepository) UpdateScore(ctx context.Context, update leaguematch.IncreaseTeamScore) error {

	exec := f.GetExecutor(ctx)
	query := `
		UPDATE score_board
		SET 
			played = played+1,
			won_sets = won_sets + $1,
			lost_sets = lost_sets + $2,
			won_games = won_games + $3,
			lost_games = lost_games + $4,
			score = score + $5,
			won = won + CASE WHEN $6 THEN 1 ELSE 0 END,
			lost = lost + CASE WHEN $6 THEN 0 ELSE 1 END
		WHERE league_id = $7
		  AND team_id = $8
	`

	result, err := exec.ExecContext(
		ctx,
		query,
		update.WonSets,
		update.LostSets,
		update.WonGames,
		update.LostGames,
		update.IncreaseScore,
		update.Won,
		update.LeagueId,
		update.TeamId,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("scoreboard not found for team %s", update.TeamId)
	}

	return nil
}

func (f *ScoreBoardRepository) InitializeScoreboard(ctx context.Context, leagueId string) error {
	exec := f.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// INSERT edilecek değerleri sağlayan SELECT sorgusu
	selectBuilder := psql.Select("league_id", "id").
		From("attendance").
		Where(squirrel.Eq{"league_id": leagueId})

	// Ana INSERT sorgusu
	insertBuilder := psql.Insert("score_board").
		Columns("league_id", "attendance_id").
		Select(selectBuilder).
		Suffix("ON CONFLICT (league_id, attendance_id) DO NOTHING")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return errors.Wrap(err, "scoreboard başlatma sorgusu oluşturulamadı")
	}

	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "scoreboard başlatılırken veritabanı hatası")
	}
	return nil
}

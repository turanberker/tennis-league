package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"fmt"
	"log"

	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/service/internal/domain/match"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
)

type MatchRepository struct {
	sqlrepository.Repository
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *MatchRepository) SaveBulkMatches(ctx context.Context, req *match.BulkInsertMatches) error {

	if len(req.Sides) == 0 {
		return nil
	}

	executor := r.GetExecutor(ctx)

	// 1. Kolon isimlerini MatchType'a göre belirle
	var side1Col, side2Col string
	if req.Type.Type == match.MatchType_SINGLE {
		side1Col = "player_1_id"
		side2Col = "player_2_id"
	} else {
		// Double veya Mix ise team kolonu kullan
		side1Col = "team_1_id"
		side2Col = "team_2_id"
	}
	// Ana sütun listesini oluştur
	columns := []string{side1Col, side2Col, "source", "match_type"}
	// Opsiyonel Lig/Turnuva ID kolonu
	if req.Type.Source == match.MatchSource_LEAGUE {
		columns = append(columns, "league_id")
	} else {
		panic("Not İmplemented Yet")
	}

	// Squirrel (sq) kullanarak bulk insert oluşturuyoruz
	builder := squirrel.Insert("tennisleague.match")
	builder = builder.PlaceholderFormat(squirrel.Dollar)
	builder = builder.Columns(columns...)

	// 2. Satırları (Values) ekle
	for _, side := range req.Sides {
		vals := []interface{}{
			side.Side1,
			side.Side2,
			req.Type.Source,
			req.Type.Type,
		}

		if req.Type.Source != match.MatchSource_FRIENDLY {
			vals = append(vals, req.Type.Id)
		}

		builder = builder.Values(vals...)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		fmt.Println(query)
		return fmt.Errorf("sorgu olusturulamadi: %w", err)
	}

	_, err = executor.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Maçlar oluşturulurken hata oluştu: %v\n", err)
		return err
	}

	return nil
}

func (r *MatchRepository) GetFixtureByLeagueId(ctx context.Context, leagueId string, filterParam *match.FixtureFilter) ([]*match.LeagueFixtureMatch, error) {
	executor := r.GetExecutor(ctx)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	// 1. SQL Builder Hazırlığı
	sqlBuilder := psql.Select(
		"m.id",
		"m.team_1_id",
		"t1.name as team_1_name",
		"m.team_1_score",
		"(m.winner_id = m.team_1_id) as team_1_is_winner",
		"m.team_2_id",
		"t2.name as team_2_name",
		"m.team_2_score",
		"(m.winner_id = m.team_2_id) as team_2_is_winner",
		"m.status",
		"m.match_date",
	).
		From("tennisleague.match m").
		Join("tennisleague.team t1 ON m.team_1_id = t1.id").
		Join("tennisleague.team t2 ON m.team_2_id = t2.id").
		Where(squirrel.Eq{"m.league_id": leagueId})

	// Dinamik filtre: TeamId null değilse sorguya ekle
	if filterParam != nil && filterParam.TeamId != nil {
		sqlBuilder = sqlBuilder.Where(squirrel.Or{
			squirrel.Eq{"m.team_1_id": *filterParam.TeamId},
			squirrel.Eq{"m.team_2_id": *filterParam.TeamId},
		})
	}

	query, args, err := sqlBuilder.OrderBy("m.match_date ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %w", err)
	}
	// 2. Metod İçi Yerel Struct (sqlscan için)
	type row struct {
		ID         string             `db:"id"`
		T1Id       string             `db:"team_1_id"`
		T1Name     string             `db:"team_1_name"`
		T1Score    *int8              `db:"team_1_score"`
		T1IsWinner *bool              `db:"team_1_is_winner"`
		T2Id       string             `db:"team_2_id"`
		T2Name     string             `db:"team_2_name"`
		T2Score    *int8              `db:"team_2_score"`
		T2IsWinner *bool              `db:"team_2_is_winner"`
		Status     match.MATCH_Status `db:"status"`
		MatchDate  *time.Time         `db:"match_date"` // Nullable tarih güvenliği
	}

	var rowsData []row
	err = sqlscan.Select(ctx, executor, &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası (fikstür): %w", err)
	}

	// 3. Domain Modeline Mapping
	var matches []*match.LeagueFixtureMatch
	for _, d := range rowsData {
		matches = append(matches, &match.LeagueFixtureMatch{
			Id: d.ID,
			Team1: match.TeamRef{
				MatchSide: match.MatchSide{
					Id:   d.T1Id,
					Name: d.T1Name,
				},
				Score:  d.T1Score,
				Winner: d.T1IsWinner,
			},
			Team2: match.TeamRef{
				MatchSide: match.MatchSide{
					Id:   d.T2Id,
					Name: d.T2Name,
				},
				Score:  d.T2Score,
				Winner: d.T2IsWinner,
			},
			Status:    d.Status,
			MatchDate: d.MatchDate,
		})
	}

	return matches, nil
}

func (r *MatchRepository) UpdateMatchDate(ctx context.Context, data match.UpdateMatchDate) error {
	executor := r.GetExecutor(ctx)
	query := "Update match set match_date=$1 where id=$2"

	result, err := executor.ExecContext(ctx, query, data.MatchDate, data.Id)
	if err != nil {
		log.Println("Maç tarihi güncellenirken hata oluştu:", err)
		return err
	}

	// Etkilenen satır sayısını kontrol et
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Nadir bir hata durumudur (sürücü kaynaklı)
	}

	if rowsAffected == 0 {
		// Kayıt bulunamadığında dönecek özel bir hata
		return fmt.Errorf("güncellenecek maç bulunamadı (ID: %s)", data.Id)
	}

	return nil
}

func (r *MatchRepository) GetMatchTeamIds(ctx context.Context, matchId string) *match.MatchTeamIds {
	executor := r.GetExecutor(ctx)
	var response match.MatchTeamIds
	query := "select league_id, team_1_id ,team_2_id ,status  from match m where id=$1"

	err := executor.QueryRowContext(ctx, query, matchId).
		Scan(&response.LeagueId, &response.Team1Id, &response.Team2Id, &response.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			// kayıt yok
			return nil
		}
		return nil
	}
	return &response
}

func (r *MatchRepository) UpdateMatchScore(ctx context.Context, macScore *match.UpdateMatchScore) error {
	executor := r.GetExecutor(ctx)
	query := "Update match set team_1_score=$1, team_2_score=$2, winner_id=$3, status=$4 where id=$5"

	_, err := executor.ExecContext(ctx, query, macScore.Team1Score, macScore.Team2Score, macScore.WinnerTeamId, match.StatusCompleted, macScore.Id)
	if err != nil {
		log.Printf("Maç Skoru güncellenirken hata oluştu:%+v", err)
		return err
	}
	return nil
}

func (r *MatchRepository) ApproveScore(ctx context.Context, source match.Match_SOURCE, matchId string) error {
	executor := r.GetExecutor(ctx)
	query := "Update match set status=$1, approve_date=current_date where id=$2 and status=$3 and source=$4"

	result, err := executor.ExecContext(
		ctx,
		query,
		match.StatusApproved, // yeni status
		matchId,
		match.StatusCompleted, // eski status
		source,
	)
	if err != nil {
		log.Printf("Maç Skoru Onaylanırken hata oluştu:%+v", err)
		return err
	}

	c, err := result.RowsAffected()
	if err != nil {
		log.Printf("Maç Skoru Onaylanırken hata oluştu:%+v", err)
		return err
	}
	if c == 0 {
		return errors.New("Onaylanacak Maç bulunamadı")
	}
	return nil
}

func (r *MatchRepository) GetMatchType(ctx context.Context, matchId string) (*match.Match_TYPE, error) {
	executor := r.GetExecutor(ctx)
	var response match.Match_TYPE
	query := "select match_type  from match m where id=$1"

	err := executor.QueryRowContext(ctx, query, matchId).
		Scan(&response)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *MatchRepository) GetPlayersIdsAndWinnerStatus(ctx context.Context, matchID string) ([]match.MatchParticipant, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	isWinnerCase := squirrel.Case().
		When("t.id IS NOT NULL", "t.id = m.winner_id").
		Else("singles.player_id = m.winner_id")
	caseSql, _, err := isWinnerCase.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build case clause: %w", err)
	}
	query, args, err := psql.
		Select(
			"COALESCE(tp.player_id, singles.player_id) AS player_id",
			fmt.Sprintf("(%s) AS is_winner", caseSql),
		).
		From("match m").
		LeftJoin("team t ON t.id IN (m.team_1_id, m.team_2_id)").
		LeftJoin("team_player tp ON tp.team_id = t.id").
		// LATERAL bloğunu ham SQL (Raw) olarak LeftJoin içine gömüyoruz
		LeftJoin("LATERAL (VALUES (m.player_1_id), (m.player_2_id)) AS singles(player_id) ON m.player_1_id IS NOT NULL").
		Where(squirrel.Eq{"m.id": matchID}).
		// Dolu satırları filtreleyen WHERE koşulu
		Where("(tp.player_id IS NOT NULL OR singles.player_id IS NOT NULL)").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// 2. sqlscan (scany) ile <sorguyu çalıştır ve sonuçları bind et
	// r.db burada *sql.DB veya *sql.Tx olabilir, scany ikisini de destekler

	type dbRow struct {
		PlayerID string `db:"player_id"`
		IsWinner bool   `db:"is_winner"`
	}

	var rows []dbRow
	if err := sqlscan.Select(ctx, r.DB, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to select participants: %w", err)
	}

	// 3. Database row'larını service modeline (match.MatchParticipant) dönüştür
	result := make([]match.MatchParticipant, len(rows))
	for i, row := range rows {
		result[i] = match.MatchParticipant{
			PlayerID: row.PlayerID,
			IsWinner: row.IsWinner,
		}
	}

	return result, nil

}

func (r *MatchRepository) GetPlayerIncomingMatches(ctx context.Context, queryParam match.PlayerIncomingMatchesQueryParam) ([]match.PlayerIncomingMatchesResult, error) {

	psql := squirrel.StatementBuilder

	// 1. SINGLE Maçlar Sorgusu
	singleMatches := psql.Select(
		"m.id AS match_id",
		"m.match_date",
		"m.match_type",
		"m.status",
		"m.source",
		"m.league_id",
		"opp.id AS opponent_id",
		"CONCAT(opp.name, ' ', opp.surname) AS opponent_name",
	).From("tennisleague.match m").
		Join("tennisleague.player opp ON (opp.id = CASE WHEN m.player_1_id = ? THEN m.player_2_id ELSE m.player_1_id END)", queryParam.PlayerId).
		Where(squirrel.Eq{"m.match_type": match.MatchType_SINGLE}).
		Where(squirrel.Or{
			squirrel.Eq{"m.player_1_id": queryParam.PlayerId},
			squirrel.Eq{"m.player_2_id": queryParam.PlayerId},
		})

	// 2. DOUBLE/TEAM Maçlar Sorgusu
	teamMatches := psql.Select(
		"m.id AS match_id",
		"m.match_date",
		"m.match_type",
		"m.status",
		"m.source",
		"m.league_id",
		"t.id AS opponent_id",
		"t.name AS opponent_name",
	).From("tennisleague.match m").
		Join("tennisleague.team_player tp ON (tp.team_id = m.team_1_id OR tp.team_id = m.team_2_id)").
		Join("tennisleague.team t ON (t.id = CASE WHEN tp.team_id = m.team_1_id THEN m.team_2_id ELSE m.team_1_id END)").
		Where(squirrel.Eq{"tp.player_id": queryParam.PlayerId}).
		Where(squirrel.Eq{"m.match_type": []match.Match_TYPE{match.MatchType_DOUBLE, match.MatchType_TEAM}})

	// UNION ALL oluşturma
	// Squirrel doğrudan UNION ALL builder'a sahip değilse, ToSql ile birleştirilir:
	singleSql, singleArgs, _ := singleMatches.ToSql()
	teamSql, teamArgs, _ := teamMatches.ToSql()

	unionSql := fmt.Sprintf("(%s UNION ALL %s)", singleSql, teamSql)
	allArgs := append(singleArgs, teamArgs...)

	// 3. Ana Sorgu (League Join ve Filtreler)
	finalQueryBuilder := psql.Select(
		"pm.match_id",
		"pm.match_date",
		"pm.match_type",
		"pm.source",
		"l.id as league_id",
		"l.name as league_name",
		"pm.opponent_id",
		"pm.opponent_name",
	).From("(" + unionSql + ") AS pm"). // Parantez içine dikkat
						LeftJoin("tennisleague.league l ON l.id = pm.league_id").
						Where("pm.match_date IS NOT NULL").
						Where(squirrel.Eq{"pm.status": match.StatusPending}).
						OrderBy("pm.match_date ASC").
						Limit(uint64(queryParam.Limit)).
						PlaceholderFormat(squirrel.Dollar)

	query, args, err := finalQueryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	allArgs = append(allArgs, args...) // UNION sorgusunun argümanları + final sorgunun argümanları
	type dbRow struct {
		MatchId      string             `db:"match_id"`
		MatchDate    *time.Time         `db:"match_date"`
		MatchType    match.Match_TYPE   `db:"match_type"`
		Source       match.Match_SOURCE `db:"source"`
		LeagueId     *string            `db:"league_id"`
		LeagueName   *string            `db:"league_name"`
		OppenentId   string             `db:"opponent_id"`
		OppenentName string             `db:"opponent_name"`
	}

	var rows []dbRow
	if err := sqlscan.Select(ctx, r.DB, &rows, query, allArgs...); err != nil {
		return nil, fmt.Errorf("failed to select participants: %w", err)
	}

	// 3. Database row'larını service modeline (match.MatchParticipant) dönüştür
	result := make([]match.PlayerIncomingMatchesResult, len(rows))
	for i, row := range rows {
		result[i] = match.PlayerIncomingMatchesResult{
			MatchId:      row.MatchId,
			MatchDate:    row.MatchDate,
			MatchType:    row.MatchType,
			Source:       row.Source,
			LeagueId:     row.LeagueId,
			LeagueName:   row.LeagueName,
			OppenentId:   row.OppenentId,
			OppenentName: row.OppenentName,
		}
	}

	return result, nil

}

func (r *MatchRepository) GetMatchInfo(ctx context.Context, matchId string) (*match.MatchInfo, error) {
	executor := r.GetExecutor(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.
		Select(
			"match_date",
			"team_1_id",
			"team_2_id",
			"player_1_id",
			"player_2_id",
			"source",
			"m.league_id",
			"m.status",
			"tournament_id",
			"match_type",
			"t1.name as team1_name",
			"t2.name as team2_name",
			"p1.name as player1_name",
			"p1.surname as player1_surname",
			"p2.name as player2_name",
			"p2.surname as player2_surname",
		).
		From("match m").
		LeftJoin("team t1 on t1.id = m.team_1_id").
		LeftJoin("team t2 on t2.id = m.team_2_id").
		LeftJoin("player p1 ON player_1_id = p1.id").
		LeftJoin("player p2 ON player_2_id = p2.id").
		Where(squirrel.Eq{"m.id": matchId}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var dbRow struct {
		MatchDate      *time.Time         `db:"match_date"`
		Team1Id        *string            `db:"team_1_id"`
		Team2Id        *string            `db:"team_2_id"`
		Player1Id      *string            `db:"player_1_id"`
		Player2Id      *string            `db:"player_2_id"`
		Team1Name      *string            `db:"team1_name"`
		Team2Name      *string            `db:"team2_name"`
		Player1Name    *string            `db:"player1_name"`
		Player2Name    *string            `db:"player2_name"`
		Player1Surname *string            `db:"player1_surname"`
		Player2Surname *string            `db:"player2_surname"`
		Source         match.Match_SOURCE `db:"source"`
		Status         match.MATCH_Status `db:"status"`
		Match_Type     match.Match_TYPE   `db:"match_type"`
		LeagueId       *string            `db:"league_id"`
		TournamentId   *string            `db:"tournament_id"`
	}

	matchInfo := &match.MatchInfo{}

	if err := sqlscan.Get(ctx, executor, &dbRow, query, args...); err != nil {
		return nil, fmt.Errorf("Maç Bilgisi Getirilemedi (id: %s): %w", matchId, err)
	}
	matchInfo.MatchDate = dbRow.MatchDate
	matchInfo.MatchType = dbRow.Match_Type
	matchInfo.Source = dbRow.Source
	matchInfo.Status = dbRow.Status
	switch dbRow.Source {
	case match.MatchSource_LEAGUE:
		if dbRow.LeagueId != nil {
			matchInfo.SourceId = dbRow.LeagueId
		}
	case match.MatchSource_PLAYOFF:
		if dbRow.TournamentId != nil {
			matchInfo.SourceId = dbRow.TournamentId
		}
	default:
		// Opsiyonel: Bilinmeyen bir kaynak gelirse loglayabilir
		// veya default bir değer atayabilirsiniz.
		matchInfo.SourceId = nil
	}
	if dbRow.Team1Id != nil && dbRow.Team2Id != nil {
		matchInfo.Side1.Id = *dbRow.Team1Id
		matchInfo.Side1.Name = *dbRow.Team1Name
		matchInfo.Side2.Id = *dbRow.Team2Id
		matchInfo.Side2.Name = *dbRow.Team2Name

	} else if dbRow.Player1Id != nil && dbRow.Player2Id != nil {
		matchInfo.Side1.Id = *dbRow.Player1Id
		matchInfo.Side1.Name = fmt.Sprintf("%s %s", *dbRow.Player1Name, *dbRow.Player1Surname)
		matchInfo.Side2.Id = *dbRow.Player2Id
		matchInfo.Side2.Name = fmt.Sprintf("%s %s", *dbRow.Player2Name, *dbRow.Player2Surname)
	} else {
		return nil, fmt.Errorf("match sides not found for match id: %s", matchId)
	}

	return matchInfo, nil

}

func (r *MatchRepository) CheckIfPlayerPlayedInMatch(ctx context.Context, matchID string, playerID string) (bool, error) {
	// Squirrel Builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Sorgu oluşturma
	query := psql.Select("COUNT(1)").
		From("\"match\" m").
		LeftJoin("team_player tp ON (tp.team_id = m.team_1_id OR tp.team_id = m.team_2_id)").
		Where(squirrel.Eq{"m.id": matchID}).
		Where(squirrel.Or{
			squirrel.Eq{"m.player_1_id": playerID},
			squirrel.Eq{"m.player_2_id": playerID},
			squirrel.Eq{"tp.player_id": playerID},
		})

	// SQL ve Argümanları al
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("query building error: %w", err)
	}

	var count int
	err = r.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("execution error: %w", err)
	}

	return count > 0, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
)

type PlayerRepository struct {
	BaseRepository
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		BaseRepository: BaseRepository{db: db}, // Base'deki db'yi dolduruyoruz
	}
}

func (r *PlayerRepository) GetById(ctx context.Context, id int64) (*player.Player, error) {
	player := &player.Player{}
	query := `SELECT id,  name, surname, sex,user_id FROM player WHERE id=$1`
	err := r.GetExecutor(ctx).QueryRowContext(ctx, query, id).
		Scan(&player.ID, &player.Name, &player.Surname, &player.Sex, &player.UserId)
	if err != nil {
		log.Println("Playerı maplerken hata oluştu:", err)
		return nil, err
	}
	return player, nil
}

func (r *PlayerRepository) List(ctx context.Context, queryParams player.ListQueryParameters) ([]*player.Player, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlBuilder := psql.Select("id", "name", "surname", "sex", "user_id", "single_point", "double_point").From("player")

	if queryParams.Name != nil {
		sqlBuilder = sqlBuilder.Where("name ILIKE ?", "%"+*queryParams.Name+"%")
	}

	if queryParams.Sex != nil {
		sqlBuilder = sqlBuilder.Where(squirrel.Eq{"sex": *queryParams.Sex})
	}
	if queryParams.HasUser != nil {
		if *queryParams.HasUser {
			sqlBuilder = sqlBuilder.Where("user_id IS NOT NULL")
		} else {
			sqlBuilder = sqlBuilder.Where("user_id IS NULL")
		}
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sorgu oluşturulamadı: %v", err)
	}

	type row struct {
		ID           string  `db:"id"`
		Name         string  `db:"name"`
		Surname      string  `db:"surname"`
		Sex          string  `db:"sex"`
		UserID       *string `db:"user_id"` // Null gelebileceği için pointer kullanmak güvenlidir
		SinglePoints int     `db:"single_point"`
		DoublePoints int     `db:"double_point"`
	}
	var rowsData []row
	err = sqlscan.Select(ctx, r.GetExecutor(ctx), &rowsData, query, args...)
	if err != nil {
		return nil, fmt.Errorf("veritabanı hatası: %v", err)
	}

	// Gelen ham veriyi kendi Player modelinize dönüştürme (Mapping)
	var players []*player.Player
	for _, d := range rowsData {
		players = append(players, &player.Player{
			ID:           d.ID,
			Name:         d.Name,
			Surname:      d.Surname,
			Sex:          player.Sex(d.Sex),
			UserId:       d.UserID,
			SinglePoints: d.SinglePoints,
			DoublePoints: d.DoublePoints,
		})
	}

	return players, nil
}

func (r *PlayerRepository) Save(ctx context.Context, persistPlayer *player.PersistPlayer) (*string, error) {
	exec := r.GetExecutor(ctx)

	query := `INSERT INTO player ( name, surname,sex, user_id) VALUES ($1, $2, $3,$4) RETURNING id`
	var id string
	err := exec.QueryRowContext(ctx, query, persistPlayer.Name, persistPlayer.Surname, persistPlayer.Sex, persistPlayer.UserId).
		Scan(&id)
	if err != nil {
		log.Println("Player insert ederken hata oluştu:", err)
		return nil, err
	}
	return &id, nil
}

func (r *PlayerRepository) AssignToUser(ctx context.Context, playerId string, userId string) error {

	exec := r.GetExecutor(ctx)
	query := `UPDATE player SET user_id = $1 WHERE id = $2`

	result, err := exec.ExecContext(ctx, query, userId, playerId)
	if err != nil {
		return fmt.Errorf("Oyuncuya kullanıcı ataması başarısız:")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return fmt.Errorf("Güncellenecek oyuncu bulunamadı: %w", err)
	}

	return nil

}

func (r PlayerRepository) DecreaseDoublePoint(ctx context.Context, playerId string, change int) (int, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.Update("player").
		Set("double_point", squirrel.Expr("GREATEST(0, double_point - ?)", change)).
		Where(squirrel.Eq{"id": playerId}).
		Suffix("RETURNING double_point").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("sorgu inşa hatası: %w", err)
	}

	exec := r.GetExecutor(ctx)

	var newPoint int
	err = exec.QueryRowContext(ctx, query, args...).Scan(&newPoint)
	if err != nil {
		// Oyuncu bulunamadıysa veya başka bir DB hatası varsa
		return 0, fmt.Errorf("oyuncu puanı düşürülemedi (ID: %s): %w", playerId, err)
	}

	return newPoint, nil
}

func (r PlayerRepository) IncreaseDoublePoint(ctx context.Context, playerId string, change int) (int, error) {
	// Squirrel builder'ı başlatalım (PostgreSQL placeholder'ları için Dollar formatı)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Sorguyu inşa edelim
	query, args, err := psql.Update("player").
		Set("double_point", squirrel.Expr("double_point + ?", change)).
		Where(squirrel.Eq{"id": playerId}).
		Suffix("RETURNING double_point").
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("sorgu oluşturulamadı: %w", err)
	}

	exec := r.GetExecutor(ctx)

	var newDoublePoint int
	// QueryRowContext kullanarak RETURNING ile gelen değeri okuyoruz
	err = exec.QueryRowContext(ctx, query, args...).Scan(&newDoublePoint)
	if err != nil {
		// Eğer hiçbir satır güncellenmediyse (id yanlışsa) sql.ErrNoRows döner
		return 0, fmt.Errorf("oyuncu puanı güncellenemedi veya oyuncu bulunamadı: %w", err)
	}

	return newDoublePoint, nil
}

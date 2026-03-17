package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
	user "github.com/turanberker/tennis-league-service/internal/domain/user"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{BaseRepository{db: db}}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.LoginUserCheck, error) {
	exec := r.GetExecutor(ctx)

	u := &user.LoginUserCheck{}
	query := `SELECT u.id, u.email,  u.name,u.surname, u.password_hash, u.role, p.id as player_id FROM "user" u 
		left join player p on p.user_id=u.id WHERE email=$1`
	err := exec.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.PasswordHash, &u.Role, &u.PlayerId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, u *user.PersistUser) (string, error) {
	exec := r.GetExecutor(ctx)

	var userId string
	query := `INSERT INTO "user" (email,  name, surname,password_hash, role) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := exec.QueryRowContext(ctx, query, u.Email, u.Name, u.Surname, u.PasswordHash, u.Role).Scan(&userId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505: unique_violation (Benzersizlik kısıtlaması ihlali)
			if pqErr.Code == "23505" {
				return "", user.USER_EXISTS_ERROR
			}
		}
		return "", err
	}

	return userId, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	exec := r.GetExecutor(ctx)
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "user"
	 WHERE email=$1)`
	err := exec.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	exec := r.GetExecutor(ctx)
	query := `SELECT u.id, u.email,  u.name,u.surname,  u.role,u.approved, p.id as player_id FROM "user" u 
		left join player p on p.user_id=u.id`
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		log.Println("Player listesi çekerken hata oluştu:", err)
		return nil, err
	}

	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		user := &user.User{}
		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.Role,
			&user.Approved,
			&user.PlayerId)

		if err != nil {
			log.Println("Playerları maplerken hata oluştu:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}

func (r *UserRepository) UpdateRoleAsCoordinator(ctx context.Context, userId string) error {
	exec := r.GetExecutor(ctx)
	query := `UPDATE "user" SET role = $1 WHERE id = $2 and role = $3`

	var err error

	_, err = exec.ExecContext(ctx, query, user.RoleCoordinator, userId, user.RolePlayer)

	if err != nil {
		return fmt.Errorf("kullanıcı rolü güncellenirken hata: %w", err)
	}

	return nil
}

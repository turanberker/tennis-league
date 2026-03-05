package postgres

import (
	"context"
	"database/sql"

	"github.com/turanberker/tennis-league-service/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u := &user.User{}
	query := `SELECT u.id, u.email, u.phone, u.name,u.surname, u.password_hash, u.role, p.id as player_id FROM users u 
		left join players p on p.user_id=u.id WHERE email=$1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Phone, &u.Name, &u.Surname, &u.PasswordHash, &u.Role, &u.PlayerId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, tx *sql.Tx, u *user.PersistUser) (string, error) {

	var userId string
	query := `INSERT INTO users (email,  name, surname,password_hash, role) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := tx.QueryRowContext(ctx, query, u.Email, u.Name, u.Surname, u.PasswordHash, u.Role).Scan(&userId)
	return userId, err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

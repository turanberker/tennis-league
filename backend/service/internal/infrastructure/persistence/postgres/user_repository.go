package postgres

import (
	"context"
	"database/sql"
	"log"

	user "github.com/turanberker/tennis-league-service/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.LoginUserCheck, error) {
	u := &user.LoginUserCheck{}
	query := `SELECT u.id, u.email,  u.name,u.surname, u.password_hash, u.role, p.id as player_id FROM users u 
		left join players p on p.user_id=u.id WHERE email=$1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Name, &u.Surname, &u.PasswordHash, &u.Role, &u.PlayerId)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return false
	}
	return exists
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	
	query := `SELECT u.id, u.email,  u.name,u.surname,  u.role,u.approved, p.id as player_id FROM users u 
		left join players p on p.user_id=u.id`
	rows, err := r.db.QueryContext(ctx, query)
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

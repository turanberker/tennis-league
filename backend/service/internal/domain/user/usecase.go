package user

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Usecase struct {
	db   *sql.DB
	repo Repository
}

func NewUsecase(db *sql.DB, r Repository) *Usecase {
	return &Usecase{db: db, repo: r}
}

func (u *Usecase) Login(ctx context.Context, email, password string) (*User, error) {
	usr, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return usr, nil
}

func (u Usecase) RegisterUser(ctx context.Context, req *RegisterUserInput) (*User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	usr, err := u.repo.SaveUser(ctx, tx, &User{
		Email:        req.Email,
		Name:         req.Name,
		Surname:      req.Surname,
		PasswordHash: string(hashedPassword),
		Role:         RolePlayer,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return usr, nil
}

package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/turanberker/tennis-league-service/internal/domain/session"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Usecase struct {
	db                *sql.DB
	repo              user.Repository
	sessionRepository session.Repository
}

func NewUsecase(db *sql.DB, r user.Repository, sessionRepository session.Repository) *Usecase {
	return &Usecase{db: db, repo: r, sessionRepository: sessionRepository}
}

func (u *Usecase) Login(ctx context.Context, email, password string) (*user.LoggedInUser, error) {
	usr, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	session, err := u.sessionRepository.Start(ctx, usr.ID, string(usr.Role), usr.PlayerId)
	if err != nil {
		return nil, err
	}
	var response = &user.LoggedInUser{
		SessionId: session.SessionId,
		ID:        usr.ID,
		Name:      usr.Name,
		Surname:   usr.Surname,
		Role:      usr.Role,
		PlayerId:  usr.PlayerId,
	}
	return response, nil
}

func (u Usecase) RegisterUser(ctx context.Context, req *user.RegisterUserInput) (*user.LoggedInUser, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userId, err := u.repo.SaveUser(ctx, tx, &user.PersistUser{
		Email:        req.Email,
		Name:         req.Name,
		Surname:      req.Surname,
		PasswordHash: string(hashedPassword),
		Role:         user.RolePlayer,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	session, err := u.sessionRepository.Start(ctx, userId, string(user.RolePlayer), nil)
	if err != nil {
		return nil, err
	}
	return &user.LoggedInUser{
		Name:      req.Name,
		Surname:   req.Surname,
		ID:        userId,
		SessionId: session.SessionId,
		Role:      user.RolePlayer,
		PlayerId:  nil}, nil
}

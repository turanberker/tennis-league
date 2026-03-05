package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/turanberker/tennis-league-service/internal/domain/session"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Usecase struct {
	db                *sql.DB
	repo              Repository
	sessionRepository session.Repository
}

func NewUsecase(db *sql.DB, r Repository, sessionRepository session.Repository) *Usecase {
	return &Usecase{db: db, repo: r, sessionRepository: sessionRepository}
}

func (u *Usecase) Login(ctx context.Context, email, password string) (*LoggedInUser, error) {
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
	var response = &LoggedInUser{
		SessionId: session.SessionId,
		ID:        usr.ID,
		Name:      usr.Name,
		Surname:   usr.Surname,
		Role:      usr.Role,
		PlayerId:  usr.PlayerId,
	}
	return response, nil
}

func (u Usecase) RegisterUser(ctx context.Context, req *RegisterUserInput) (*LoggedInUser, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userId, err := u.repo.SaveUser(ctx, tx, &PersistUser{
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
	session, err := u.sessionRepository.Start(ctx, userId, string(RolePlayer), nil)
	if err != nil {
		return nil, err
	}
	return &LoggedInUser{
		Name:      req.Name,
		Surname:   req.Surname,
		ID:        userId,
		SessionId: session.SessionId,
		Role:      RolePlayer,
		PlayerId:  nil}, nil
}

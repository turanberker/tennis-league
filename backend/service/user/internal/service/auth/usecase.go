package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"tennis-league/common/security/dto"
	service "tennis-league/user-service/internal"
	"tennis-league/user-service/internal/service/session"
	userService "tennis-league/user-service/internal/service/user"

	customerror "tennis-league/common/lib/error"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Usecase struct {
	db                *sql.DB
	repo              userService.Repository
	sessionRepository session.Repository
}

func NewUsecase(db *sql.DB, r userService.Repository, sessionRepository session.Repository) *Usecase {
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

	startSessionInput := session.StartSessionInput{
		UserId:   usr.ID,
		Role:     string(usr.Role),
		PlayerId: usr.PlayerId,
	}

	session, err := u.sessionRepository.Start(ctx, &startSessionInput)
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

	userId, err := u.repo.SaveUser(ctx, &userService.PersistUser{
		Email:        req.Email,
		Name:         req.Name,
		Surname:      req.Surname,
		PasswordHash: string(hashedPassword),
		Role:         dto.RolePlayer,
	})
	if err != nil {
		if errors.Is(err, userService.USER_EXISTS_ERROR) {
			return nil, &customerror.BusinnesException{
				StatusCode: http.StatusOK,
				ErrorCode:  service.ErrCodeEmailAlreadyExists, // "EMAIL_ALREADY_EXISTS"
				Message:    "Bu e-posta adresiyle daha önce kayıt olunmuş.",
			}
		}
		return nil, err
	}

	startSessionInput := session.StartSessionInput{
		UserId:   userId,
		Role:     string(dto.RolePlayer),
		PlayerId: nil,
	}

	session, err := u.sessionRepository.Start(ctx, &startSessionInput)
	if err != nil {
		return nil, err
	}
	return &LoggedInUser{
		Name:      req.Name,
		Surname:   req.Surname,
		ID:        userId,
		SessionId: session.SessionId,
		Role:      dto.RolePlayer,
		PlayerId:  nil}, nil
}

func (u Usecase) DeleteSessionFromRedis(ctx context.Context, sessionId string) {

	u.sessionRepository.Delete(ctx, sessionId)

}

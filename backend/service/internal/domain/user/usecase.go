package user

import (
	"context"
	"net/http"

	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	tm   *database.TransactionManager
	repo Repository
}

func (u *Usecase) ChangePassword(ctx context.Context, userId string, currentPassword string, newPassword string) error {

	user, err := u.repo.FindById(ctx, userId)
	if err != nil {
		return customerror.NewInternalError(err)
	}

	if user == nil {
		return &customerror.BusinnesException{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  customerror.ErrUserNotExists,
			Message:    "Kullanıcı Bulunamadı",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return &customerror.BusinnesException{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  customerror.ErrInvalidPassword,
			Message:    "Girmiş olduğunuz şifre yanlış",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return customerror.NewInternalError(err)
	}

	err = u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.UpdatePassword(txCtx, userId, string(hashedPassword))
	})

	if err != nil {
		return customerror.NewInternalError(err)
	}
	return nil
}

func (u *Usecase) SetUserAsCoordinator(ctx context.Context, userId string) error {
	// 1. İşlemi atomik hale getirmek için Transaction başlatıyoruz
	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		// 4. Aksiyon: Rolü güncelle
		err := u.repo.UpdateRoleAsCoordinator(txCtx, userId)
		if err != nil {
			return err
		}

		return nil
	})
}

func NewUsecase(
	tm *database.TransactionManager,
	r Repository) *Usecase {
	return &Usecase{repo: r, tm: tm}
}

func (u *Usecase) GetAll(ctx context.Context) ([]*User, error) {
	return u.repo.List(ctx)
}

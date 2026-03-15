package user

import (
	"context"

	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

type Usecase struct {
	tm   *database.TransactionManager
	repo Repository
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

func NewUsecase(r Repository,
	tm *database.TransactionManager) *Usecase {
	return &Usecase{repo: r, tm: tm}
}

func (u *Usecase) GetAll(ctx context.Context) ([]*User, error) {
	return u.repo.List(ctx)
}

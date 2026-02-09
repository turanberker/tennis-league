package league

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNameFieldRequired = errors.New("Name can not be null or empty string")
var ErrNameLenghtError = errors.New("Name size must between 5 and 75 characters")

type Usecase struct {
	db   *sql.DB
	repo Repository
}

func NewUsecase(db *sql.DB, r Repository) *Usecase {
	return &Usecase{db: db, repo: r}
}

func (u *Usecase) GetById(ctx context.Context, id int64) (*League, error) {

	return u.repo.GetById(ctx, id)
}

func (u *Usecase) GetAll(ctx context.Context, name string) ([]*League, error) {

	return u.repo.GetAll(ctx, name)
}

func (u *Usecase) Save(ctx context.Context, persistLeague *PersistLeague) (int64, error) {
	if persistLeague.Name == "" {
		return 0, ErrNameFieldRequired
	}

	if len(persistLeague.Name) < 5 || len(persistLeague.Name) > 75 {
		return 0, ErrNameLenghtError
	}

	return u.repo.Save(ctx, persistLeague)
}

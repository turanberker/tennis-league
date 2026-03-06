package user

import (
	"context"
	"database/sql"
)

type Usecase struct {
	db   *sql.DB
	repo Repository
}

func NewUsecase(db *sql.DB, r Repository) *Usecase {
	return &Usecase{db: db, repo: r}
}

func (u *Usecase) GetAll(ctx context.Context) ([]*User,error){
	return u.repo.List(ctx)
}

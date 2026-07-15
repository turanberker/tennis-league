package matchrequest

import (
	"context"
	"tennis-league/common/lib/database"
)

type UseCase struct {
	tm         *database.TransactionManager
	repository Repository
}

func NewMatchRequestUseCase(tm *database.TransactionManager, repository Repository) *UseCase {
	return &UseCase{tm: tm, repository: repository}
}

func (m *UseCase) NewRequest(ctx context.Context, request RequestDto) (string, error) {
	var id string

	err := m.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		requestId, err := m.repository.Create(txCtx, request)
		id = requestId
		return err
	})

	if err != nil {
		return id, err
	}
	return id, nil
}

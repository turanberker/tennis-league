package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sqlrepository "tennis-league/common/lib/repository/sql"
	"tennis-league/common/security/dto"
)

type UserRepository struct {
	sqlrepository.Repository
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Repository: *sqlrepository.NewRepository(db)}
}

func (r *UserRepository) UpdateRoleAsCoordinator(ctx context.Context, userId string) error {
	exec := r.GetExecutor(ctx)
	query := `UPDATE "user" SET role = $1 WHERE id = $2 and role = $3`

	var err error

	_, err = exec.ExecContext(ctx, query, dto.RoleCoordinator, userId, dto.RolePlayer)

	if err != nil {
		return fmt.Errorf("kullanıcı rolü güncellenirken hata: %w", err)
	}

	return nil
}

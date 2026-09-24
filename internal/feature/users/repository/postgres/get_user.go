package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
	core_postgres_pool "github.com/DenisAstrakhan/api-server/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		WHERE id=$1;
	`
	row := r.pool.QueryRow(timectx, query, id)
	var userModel UserModels
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
	userDpmain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)
	return userDpmain, nil
}

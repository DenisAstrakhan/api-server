package user_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
	core_postgres_pool "github.com/DenisAstrakhan/api-server/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(ctx context.Context, id int, user domain.User) (domain.User, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
	UPDATE todoapp.users
	SET
		full_name=$1,
		phone_number=$2,
		version=version+1
	WHERE ID=$3 AND version=$4
	RETURNING id, version, full_name, phone_number;
	`
	row := r.pool.QueryRow(timectx, query, user.FullName, user.PhoneNumber, id, user.Version)
	var userModel UserModels
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently accessed: %w", id, core_errors.ErrConflict)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
	userDpmain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)
	return userDpmain, nil
}

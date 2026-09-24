package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
INSERT INTO todoapp.users (full_name,phone_number)
VALUES ($1,$2)
RETURNING id, version, full_name, phone_number;
`
	row := r.pool.QueryRow(timectx, query, user.FullName, user.PhoneNumber)
	var userModel UserModels
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
	userDpmain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)
	return userDpmain, nil
}

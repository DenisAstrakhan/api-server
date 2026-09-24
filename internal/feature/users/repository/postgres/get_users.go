package user_postgres_repository

import (
	"context"
	"fmt"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
)

func (r *UsersRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	timectx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()
	query := `
		SELECT id, version, full_name, phone_number
		FROM todoapp.users
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
	`
	rows, err := r.pool.Query(timectx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var userModels []UserModels
	for rows.Next() {
		var userModel UserModels
		err := rows.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}
	userDomains := userDomainsFromModels(userModels)
	return userDomains, nil
}

package user_postgres_repository

import "github.com/DenisAstrakhan/api-server/internal/core/domain"

type UserModels struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []UserModels) []domain.User {
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = domain.NewUser(user.ID, user.Version, user.FullName, user.PhoneNumber)
	}
	return domainUsers
}

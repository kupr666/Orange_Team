package users_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type UserModel struct {
	ID               uuid.UUID
	Version          int
	Email            string
	FullName         string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	UserWorkoutScore int
	Sex              *string
	WeightGrams      *int
	BirthDate        *time.Time
	HeightCM         *int
}

func (m *UserModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.Email,
		&m.FullName,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.UserWorkoutScore,
		&m.Sex,
		&m.WeightGrams,
		&m.BirthDate,
		&m.HeightCM,
	)
}

func domainFromModel(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.Email,
		model.FullName,
		model.CreatedAt,
		model.UpdatedAt,
		model.UserWorkoutScore,
		model.Sex,
		model.WeightGrams,
		model.BirthDate,
		model.HeightCM,
	)
}

// func domainsFromModels(models []UserModel) []domain.User {
// 	domains := make([]domain.User, len(models))
// 	for i, model := range models {
// 		domains[i] = domainFromModel(model)
// 	}
// 	return domains
// }

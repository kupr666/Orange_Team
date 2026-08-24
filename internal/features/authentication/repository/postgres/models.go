package authentication_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type UserModel struct {
	id               uuid.UUID
	version          int64
	mail             string
	fullName         string
	createdAt        time.Time
	updatedAt        *time.Time
	userWorkoutScore int
	sex              string
	weightGrams      *int
	birthDate        *time.Time
	heightCm         *int
}

func domainUserFromModel(model UserModel) domain.User {
	return domain.User{
		ID:               model.id,
		Version:          model.version,
		Email:            model.mail,
		FullName:         model.fullName,
		CreatedAt:        model.createdAt,
		UpdatedAt:        model.updatedAt,
		UserWorkoutScore: model.userWorkoutScore,
		Sex:              model.sex,
		WeightGrams:      model.weightGrams,
		BirthDate:        model.birthDate,
		HeightCM:         model.heightCm,
	}
}

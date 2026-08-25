package auth_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type UserResponse struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	FullName         string     `json:"full_name"`
	ProfileCompleted bool       `json:"profile_completed"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	UserWorkoutScore int        `json:"user_workout_score"`
	Sex              *string    `json:"sex,omitempty"`
	WeightGrams      *int       `json:"weight_grams,omitempty"`
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	HeightCM         *int       `json:"height_cm,omitempty"`
}

// userDTOFromDomain преобразует доменную сущность пользователя в DTO.
func userDTOFromDomain(user domain.User) UserResponse {
	return UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		FullName:         user.FullName,
		ProfileCompleted: user.ProfileCompleted(),
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		UserWorkoutScore: user.UserWorkoutScore,
		Sex:              user.Sex,
		WeightGrams:      user.WeightGrams,
		BirthDate:        user.BirthDate,
		HeightCM:         user.HeightCM,
	}
}

// userDTOsFromDomains преобразует слайс доменных пользователей в слайс DTO.
func userDTOsFromDomains(users []domain.User) []UserResponse {
	dtos := make([]UserResponse, len(users))
	for i, u := range users {
		dtos[i] = userDTOFromDomain(u)
	}
	return dtos
}

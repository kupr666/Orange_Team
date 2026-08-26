package users_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type UserDTOResponse struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	FullName         string     `json:"full_name"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	UserWorkoutScore int        `json:"user_workout_score"`
	Sex              *string    `json:"sex,omitempty"`
	WeightGrams      *int       `json:"weight_grams,omitempty"`
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	HeightCM         *int       `json:"height_cm,omitempty"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:               user.ID,
		Email:            user.Email,
		FullName:         user.FullName,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		UserWorkoutScore: user.UserWorkoutScore,
		Sex:              user.Sex,
		WeightGrams:      user.WeightGrams,
		BirthDate:        user.BirthDate,
		HeightCM:         user.HeightCM,
	}
}

func userDTOsFromDomains(users []domain.User) []UserDTOResponse {
	userDTO := make([]UserDTOResponse, len(users))
	for i, user := range users {
		userDTO[i] = userDTOFromDomain(user)
	}
	return userDTO
}

package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
)

func (s *WorkoutsService) CreateWorkout(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Workout, error) {
	if err := validateCreateWorkout(userID); err != nil {
		return domain.Workout{}, err
	}

	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"get user: %w",
			err,
		)
	}

	coeff := calculatePersonalScoreCoefficient(user)

	workout, err := s.workoutsRepository.CreateWorkout(ctx, userID, coeff)

	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"create workout in repository: %w",
			err,
		)
	}

	return workout, nil
}

func validateCreateWorkout(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf(
			"user ID is empty: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func calculatePersonalScoreCoefficient(user domain.User) int {
	// Если профиль не заполнен — возвращаем 1
	if user.WeightGrams == nil || user.HeightCM == nil {
		return 1
	}

	weightKg := float64(*user.WeightGrams) / 1000.0
	heightM := float64(*user.HeightCM) / 100.0

	bmi := weightKg / (heightM * heightM)

	coeff := 1.0 + (bmi-18.0)*9.0/17.0
	if coeff < 1.0 {
		coeff = 1.0
	}
	if coeff > 10.0 {
		coeff = 10.0
	}

	if user.Sex != nil && *user.Sex == domain.SexMale {
		coeff += 0.5
	}
	if coeff > 10.0 {
		coeff = 10.0
	}

	return int(coeff + 0.5)
}

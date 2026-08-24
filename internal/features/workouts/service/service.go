package workouts_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutsService struct {
	workoutsRepository WorkoutsRepository
}

type WorkoutsRepository interface {
	CreateWorkout(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.Workout, error)

	GetWorkouts(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Workout, error)

	GetWorkout(
		ctx context.Context,
		workoutID uuid.UUID,
	) (domain.Workout, error)

	DeleteWorkout(
		ctx context.Context,
		workoutID uuid.UUID,
	) error

	// PatchWorkout(
	// args
	// ) (retun values)

	// DeletetWorkout(
	// args
	// ) (retun values)
}

func NewWorkoutsService(repo WorkoutsRepository) *WorkoutsService {
	return &WorkoutsService{
		workoutsRepository: repo,
	}
}

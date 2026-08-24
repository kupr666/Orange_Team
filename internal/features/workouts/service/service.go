package workouts_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutsService struct {
	workoutsRepository WorkoutsRepository

	// Подключение репозитория workoutExercisesRepository

	// workoutExercisesRepository WorkoutExercisesRepository
}

type WorkoutsRepository interface {
	// PostWorkout(
	// args
	// ) (retun values)

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

	PatchWorkout(
		ctx context.Context,
		workout domain.Workout,
	) (domain.Workout, error)
}

func NewWorkoutsService(repo WorkoutsRepository) *WorkoutsService {
	return &WorkoutsService{
		workoutsRepository: repo,
	}
}

package workouts_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutExerciseReader interface {
	GetWorkoutExercises(
		ctx context.Context,
		workoutID uuid.UUID,
	) ([]domain.WorkoutExercise, error)
}

type UserRepository interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	UpdateUserWorkoutScore(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type WorkoutsService struct {
	workoutsRepository    WorkoutsRepository
	workoutExerciseReader WorkoutExerciseReader
	userRepository        UserRepository
}

type WorkoutsRepository interface {
	CreateWorkout(
		ctx context.Context,
		userID uuid.UUID,
		personalScoreCoefficient int,
	) (domain.Workout, error)

	GetWorkouts(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Workout, error)

	GetWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) (domain.Workout, error)

	DeleteWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) error

	PatchWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workout domain.Workout,
	) (domain.Workout, error)
}

func NewWorkoutsService(
	workoutsRepository WorkoutsRepository,
	workoutExerciseReader WorkoutExerciseReader,
	userRepository UserRepository,
) *WorkoutsService {
	return &WorkoutsService{
		workoutsRepository:    workoutsRepository,
		workoutExerciseReader: workoutExerciseReader,
		userRepository:        userRepository,
	}
}

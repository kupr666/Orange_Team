package workout_exercises_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutExercisesRepository interface {
	CreateWorkoutExercise(
		ctx context.Context,
		workoutExercise domain.WorkoutExercise,
	) (domain.WorkoutExercise, error)

	GetWorkoutExercises(
		ctx context.Context,
		workoutID uuid.UUID,
	) ([]domain.WorkoutExercise, error)

	GetWorkoutExercise(
		ctx context.Context,
		workoutID uuid.UUID,
		workoutExerciseID uuid.UUID,
	) (domain.WorkoutExercise, error)

	PatchWorkoutExercise(
		ctx context.Context,
		workoutExercise domain.WorkoutExercise,
	) (domain.WorkoutExercise, error)

	DeleteWorkoutExercise(
		ctx context.Context,
		workoutID uuid.UUID,
		workoutExerciseID uuid.UUID,
	) error

	GetExercise(
		ctx context.Context,
		exerciseID uuid.UUID,
	) (domain.Exercise, error)

	GetWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) (domain.Workout, error)

	UpdateWorkoutScore(
		ctx context.Context,
		workoutID uuid.UUID,
		score int,
	) error

	GetPersonalScoreCoefficient(
		ctx context.Context,
		workoutID uuid.UUID,
	) (int, error)
}

type WorkoutExercisesService struct {
	repository WorkoutExercisesRepository
	now        func() time.Time
}

func NewWorkoutExercisesService(
	repository WorkoutExercisesRepository,
) *WorkoutExercisesService {
	return &WorkoutExercisesService{
		repository: repository,
		now:        time.Now,
	}
}

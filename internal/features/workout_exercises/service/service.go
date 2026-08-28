package workout_exercises_service

import (
	"context"

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
}

type ExerciseRepository interface {
	GetExercise(
		ctx context.Context,
		exerciseID uuid.UUID,
	) (domain.Exercise, error)
}

type WorkoutRepository interface {
	GetWorkout(
		ctx context.Context,
		workoutID uuid.UUID,
	) (domain.Workout, error)
}

type WorkoutScoreUpdater interface {
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
	workoutExercisesRepository WorkoutExercisesRepository
	exerciseRepository         ExerciseRepository
	workoutRepository          WorkoutRepository
	workoutUpdater             WorkoutScoreUpdater
}

func NewWorkoutExercisesService(
	workoutExercisesRepository WorkoutExercisesRepository,
	exercisesRepository ExerciseRepository,
	workoutRepository WorkoutRepository,
	workoutUpdater WorkoutScoreUpdater,
) *WorkoutExercisesService {
	return &WorkoutExercisesService{
		workoutExercisesRepository: workoutExercisesRepository,
		exerciseRepository:         exercisesRepository,
		workoutRepository:          workoutRepository,
		workoutUpdater:             workoutUpdater,
	}
}

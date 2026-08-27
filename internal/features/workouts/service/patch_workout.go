package workouts_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *WorkoutsService) PatchWorkout(
	ctx context.Context,
	userID uuid.UUID,
	workoutID uuid.UUID,
	patch domain.WorkoutPatch,
) (domain.Workout, error) {
	workout, err := s.workoutsRepository.GetWorkout(ctx, userID, workoutID)
	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"get workout from repository: %w",
			err,
		)
	}

	if err := workout.ApplyPatch(patch); err != nil {
		return domain.Workout{}, fmt.Errorf(
			"apply workout patch: %w",
			err,
		)
	}

	// Если статус стал "completed" — проверяем, что все упражнения выполнены.
	// Для этого нужен репозиторий workoutExercisesRepository, который можно добавить в сервис.
	// Пока оставляем заглушку, чтобы код собирался.

	// if workout.Status == "completed" {
	// 	workoutExercises, err := s.workoutExercisesRepository.GetWorkout(ctx, workoutID)
	// 	if err != nil {
	// 		return domain.Workout{}, fmt.Errorf(
	// 			"get exercises for workout: %w",
	// 			err,
	// 		)
	// 	}
	// 	for _, workoutExercise := range workoutExercises {
	// 		if !workoutExercise.Completed {
	// 			return domain.Workout{}, fmt.Errorf(
	// 				"cannot complete workout: exercise %s not completed",
	// 				workoutExercise.ID)
	// 		}
	// 	}
	// }

	updatedWorkout, err := s.workoutsRepository.PatchWorkout(ctx, userID, workout)
	if err != nil {
		return domain.Workout{}, fmt.Errorf(
			"update workout in repository: %w",
			err,
		)
	}

	return updatedWorkout, nil
}

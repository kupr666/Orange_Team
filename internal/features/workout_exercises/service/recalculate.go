package workout_exercises_service

import (
	"context"

	"github.com/google/uuid"
)

// recalculateScore — временная заглушка для пересчёта workout_score.
// В будущем здесь будет полноценная логика с учётом difficulty и personal_score_coefficient.
func (s *WorkoutExercisesService) recalculateScore(
	ctx context.Context,
	workoutID uuid.UUID,
) error {
	// Пока ничего не делаем.
	// TODO: реализовать пересчёт с учётом difficulty и personal_score_coefficient
	return nil
}

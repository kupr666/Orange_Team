package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

// PatchWorkoutExercise обновляет упражнение в тренировке.
// Временно заглушка, так как требуется метод GetByID в репозитории.
func (s *WorkoutExercisesService) PatchWorkoutExercise(
	ctx context.Context,
	userID uuid.UUID,
	workoutExerciseID uuid.UUID,
	patch domain.WorkoutExercisePatch,
) (domain.WorkoutExercise, error) {
	// TODO: реализовать после добавления GetByID в репозиторий
	// 1. Получить существующее упражнение по ID
	// 2. Применить патч (вызвать ApplyPatch)
	// 3. Сохранить обновлённое упражнение
	// 4. Пересчитать workout_score
	return domain.WorkoutExercise{}, fmt.Errorf(
		"not implemented: need GetByID in repository",
	)
}

package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	maximumWorkoutScore       = 10_000
	scoreSaturationScale      = 10_000
	maximumScaledLoad         = 10_000_000
	minimumScoreCoefficient   = 1
	maximumScoreCoefficient   = 10
	minimumExerciseDifficulty = 1
	maximumExerciseDifficulty = 10
)

func (s *WorkoutExercisesService) recalculateScore(
	ctx context.Context,
	workoutID uuid.UUID,
) error {
	exercisesWithDiff, err := s.workoutExercisesRepository.GetWorkoutExercisesWithDifficulty(ctx, workoutID)
	if err != nil {
		return fmt.Errorf(
			"get workout exercises with difficulty: %w",
			err,
		)
	}

	coefficient, err := s.workoutUpdater.GetPersonalScoreCoefficient(ctx, workoutID)
	if err != nil {
		return fmt.Errorf("get personal score coefficient: %w", err)
	}
	if coefficient < minimumScoreCoefficient || coefficient > maximumScoreCoefficient {
		return fmt.Errorf(
			"personal score coefficient must be between %d and %d, got %d",
			minimumScoreCoefficient,
			maximumScoreCoefficient,
			coefficient,
		)
	}

	// The curve uses K = coefficient / scoreSaturationScale:
	// score = maxScore * (K * load) / (1 + K * load).
	// This is equivalent to the integer formula used below and makes a larger
	// personal coefficient reach saturation faster.
	maximumEffectiveLoad := divideRoundUp(maximumScaledLoad, coefficient)
	effectiveLoad := 0

	for _, ex := range exercisesWithDiff {
		if !ex.Completed {
			continue
		}
		if ex.ExerciseLoad <= 0 {
			return fmt.Errorf(
				"completed workout exercise '%s' has non-positive load: %d",
				ex.ID,
				ex.ExerciseLoad,
			)
		}
		if ex.Difficulty < minimumExerciseDifficulty || ex.Difficulty > maximumExerciseDifficulty {
			return fmt.Errorf(
				"exercise '%s' difficulty must be between %d and %d, got %d",
				ex.ExerciseID,
				minimumExerciseDifficulty,
				maximumExerciseDifficulty,
				ex.Difficulty,
			)
		}
		effectiveLoad = addLoadWithLimit(
			effectiveLoad,
			ex.ExerciseLoad,
			ex.Difficulty,
			maximumEffectiveLoad,
		)
		if effectiveLoad == maximumEffectiveLoad {
			break
		}
	}

	score := saturatedScore(effectiveLoad, coefficient)
	if err := s.workoutUpdater.UpdateWorkoutScore(ctx, workoutID, score); err != nil {
		return fmt.Errorf("update workout score: %w", err)
	}

	return nil
}

func addLoadWithLimit(current, load, difficulty, limit int) int {
	// if effectiveLoad >= maximumEffectiveLoad ||
	// workoutExercise.ExerciseLoad > (maximumEffectiveLoad - effectiveLoad) / difficulty
	if current >= limit || load > (limit-current)/difficulty {
		return limit
	}

	// effectiveLoad + workoutExercise.ExerciseLoad * difficulty
	return current + load*difficulty
}

func saturatedScore(effectiveLoad, coefficient int) int {
	scaledLoad := effectiveLoad * coefficient
	denominator := scoreSaturationScale + scaledLoad

	// Adding denominator/2 implements deterministic rounding to the nearest
	// integer. maximumScaledLoad caps the intermediate values and is the point
	// at which the curve rounds to its 10_000-point saturation value.
	score := (maximumWorkoutScore*scaledLoad + denominator/2) / denominator
	return int(score)
}

func divideRoundUp(value, divisor int) int {
	return (value + divisor - 1) / divisor
}

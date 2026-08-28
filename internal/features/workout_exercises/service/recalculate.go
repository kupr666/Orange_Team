package workout_exercises_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	maximumWorkoutScore       int64 = 10_000
	scoreSaturationScale      int64 = 100
	maximumScaledLoad         int64 = 2_000_000
	minimumScoreCoefficient         = 1
	maximumScoreCoefficient         = 10
	minimumExerciseDifficulty       = 1
	maximumExerciseDifficulty       = 10
)

func (s *WorkoutExercisesService) recalculateScore(
	ctx context.Context,
	workoutID uuid.UUID,
) error {
	workoutExercises, err := s.workoutExercisesRepository.GetWorkoutExercises(ctx, workoutID)
	if err != nil {
		return fmt.Errorf("get workout exercises for score recalculation: %w", err)
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
	maximumEffectiveLoad := divideRoundUp(maximumScaledLoad, int64(coefficient))
	effectiveLoad := int64(0)
	difficultyByExerciseID := make(map[uuid.UUID]int)

	for _, workoutExercise := range workoutExercises {
		if !workoutExercise.Completed {
			continue
		}
		if workoutExercise.ExerciseLoad <= 0 {
			return fmt.Errorf(
				"completed workout exercise '%s' has non-positive load: %d",
				workoutExercise.ID,
				workoutExercise.ExerciseLoad,
			)
		}

		difficulty, ok := difficultyByExerciseID[workoutExercise.ExerciseID]
		if !ok {
			exercise, err := s.exerciseRepository.GetExercise(ctx, workoutExercise.ExerciseID)
			if err != nil {
				return fmt.Errorf(
					"get exercise '%s' for score recalculation: %w",
					workoutExercise.ExerciseID,
					err,
				)
			}

			difficulty = exercise.Difficulty
			if difficulty < minimumExerciseDifficulty || difficulty > maximumExerciseDifficulty {
				return fmt.Errorf(
					"exercise '%s' difficulty must be between %d and %d, got %d",
					exercise.ID,
					minimumExerciseDifficulty,
					maximumExerciseDifficulty,
					difficulty,
				)
			}
			difficultyByExerciseID[workoutExercise.ExerciseID] = difficulty
		}

		effectiveLoad = addLoadWithLimit(
			effectiveLoad,
			int64(workoutExercise.ExerciseLoad),
			int64(difficulty),
			maximumEffectiveLoad,
		)
		if effectiveLoad == maximumEffectiveLoad {
			break
		}
	}

	score := saturatedScore(effectiveLoad, int64(coefficient))
	if err := s.workoutUpdater.UpdateWorkoutScore(ctx, workoutID, score); err != nil {
		return fmt.Errorf("update workout score: %w", err)
	}

	return nil
}

func addLoadWithLimit(current, load, difficulty, limit int64) int64 {
	// if effectiveLoad >= maximumEffectiveLoad ||
	// workoutExercise.ExerciseLoad > (maximumEffectiveLoad - effectiveLoad) / difficulty
	if current >= limit || load > (limit-current)/difficulty {
		return limit
	}

	// effectiveLoad + workoutExercise.ExerciseLoad * difficulty
	return current + load*difficulty
}

func saturatedScore(effectiveLoad, coefficient int64) int {
	scaledLoad := effectiveLoad * coefficient
	denominator := scoreSaturationScale + scaledLoad

	// Adding denominator/2 implements deterministic rounding to the nearest
	// integer. maximumScaledLoad caps the intermediate values and is the point
	// at which the curve rounds to its 10_000-point saturation value.
	score := (maximumWorkoutScore*scaledLoad + denominator/2) / denominator
	return int(score)
}

func divideRoundUp(value, divisor int64) int64 {
	return (value + divisor - 1) / divisor
}

package workouts_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (r *WorkoutsRepository) GetWorkouts(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Workout, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT 
		id,
        version,
        user_id,
        status,
        started_at,
        completed_at,
        created_at,
        updated_at,
        workout_score,
        intensity,
        personal_score_coefficient
	FROM app.workouts
	WHERE user_id = $1
	ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select workouts: %w", err)
	}
	defer rows.Close()

	var workoutsModels []WorkoutModel

	for rows.Next() {
		var workoutModel WorkoutModel

		// reused Scan method for each row (method is located in models.go)
		err := workoutModel.Scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan workouts: %w", err)
		}

		workoutsModels = append(workoutsModels, workoutModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	workoutsDomains := domainsFromModels(workoutsModels)

	return workoutsDomains, nil
}

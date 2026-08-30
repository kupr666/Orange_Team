package workout_exercises_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type WorkoutExerciseModel struct {
	ID           uuid.UUID
	Version      int
	WorkoutID    uuid.UUID
	ExerciseID   uuid.UUID
	Weight       *int
	Sets         *int
	Reps         *int
	Duration     *int
	Completed    bool
	ExerciseLoad int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func (m *WorkoutExerciseModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.WorkoutID,
		&m.ExerciseID,
		&m.Weight,
		&m.Sets,
		&m.Reps,
		&m.Duration,
		&m.Completed,
		&m.ExerciseLoad,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
}

func domainFromModel(model WorkoutExerciseModel) domain.WorkoutExercise {
	return domain.NewWorkoutExercise(
		model.ID,
		model.Version,
		model.WorkoutID,
		model.ExerciseID,
		model.Weight,
		model.Sets,
		model.Reps,
		model.Duration,
		model.Completed,
		model.ExerciseLoad,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

func domainFromModels(models []WorkoutExerciseModel) []domain.WorkoutExercise {
	domains := make([]domain.WorkoutExercise, len(models))

	for i, model := range models {
		domains[i] = domainFromModel(model)
	}

	return domains
}

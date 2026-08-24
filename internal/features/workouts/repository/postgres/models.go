package workouts_postgres_repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_postgres_pool "github.com/kupr666/Orange_Team/internal/core/repository/postgres/pool"
)

type WorkoutModel struct {
	ID                       uuid.UUID
	Version                  int
	UserID                   uuid.UUID
	Status                   string
	StartedAt                *time.Time
	CompletedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	WorkoutScore             int
	Intensity                *int
	PersonalScoreCoefficient int
}

type CreatedExerciseModel struct {
	ID           uuid.UUID
	Version      int
	ExerciseID   uuid.UUID
	WorkoutID    uuid.UUID
	Weight       *int
	Sets         *int
	Reps         *int
	Duration     *int
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Completed    bool
	ExerciseLoad int
}

func (m *WorkoutModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.UserID,
		&m.Status,
		&m.StartedAt,
		&m.CompletedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.WorkoutScore,
		&m.Intensity,
		&m.PersonalScoreCoefficient,
	)
}

func (m *CreatedExerciseModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.ExerciseID,
		&m.WorkoutID,
		&m.Weight,
		&m.Sets,
		&m.Reps,
		&m.Duration,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Completed,
		&m.ExerciseLoad,
	)
}

func domainFromModel(model WorkoutModel) domain.Workout {
	return domain.NewWorkout(
		model.ID,
		model.Version,
		model.UserID,
		model.Status,
		model.StartedAt,
		model.CompletedAt,
		model.CreatedAt,
		model.UpdatedAt,
		model.WorkoutScore,
		model.Intensity,
		model.PersonalScoreCoefficient,
	)
}

func createdExerciseDomainFromModel(
	model CreatedExerciseModel,
) domain.CreatedExercise {
	return domain.NewCreatedExercise(
		model.ID,
		model.Version,
		model.ExerciseID,
		model.WorkoutID,
		model.Weight,
		model.Sets,
		model.Reps,
		model.Duration,
		model.CreatedAt,
		model.UpdatedAt,
		model.Completed,
		model.ExerciseLoad,
	)
}

func domainsFromModels(models []WorkoutModel) []domain.Workout {
	domains := make([]domain.Workout, len(models))

	for i, model := range models {
		domains[i] = domainFromModel(model)
	}

	return domains
}

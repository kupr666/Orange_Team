package workouts_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutDTOResponse struct {
	ID                       uuid.UUID  `json:"id"`
	Version                  int        `json:"version"`
	UserID                   uuid.UUID  `json:"user_id"`
	Status                   string     `json:"status"`
	StartedAt                *time.Time `json:"started_at"`
	CompletedAt              *time.Time `json:"completed_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	WorkoutScore             int        `json:"workout_score"`
	Intensity                *int       `json:"intensity"`
	PersonalScoreCoefficient int        `json:"personal_score_coefficient"`
}

func workoutDTOFromDomain(workout domain.Workout) WorkoutDTOResponse {
	return WorkoutDTOResponse{
		ID:                       workout.ID,
		Version:                  workout.Version,
		UserID:                   workout.UserID,
		Status:                   workout.Status,
		StartedAt:                workout.StartedAt,
		CompletedAt:              workout.CompletedAt,
		CreatedAt:                workout.CreatedAt,
		UpdatedAt:                workout.UpdatedAt,
		WorkoutScore:             workout.WorkoutScore,
		Intensity:                workout.Intensity,
		PersonalScoreCoefficient: workout.PersonalScoreCoefficient,
	}
}

func workoutDTOsFromDomains(workouts []domain.Workout) []WorkoutDTOResponse {
	workoutDTO := make([]WorkoutDTOResponse, len(workouts))

	for i, workout := range workouts {
		workoutDTO[i] = workoutDTOFromDomain(workout)
	}

	return workoutDTO
}

type CreatedExerciseDTOResponse struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	ExerciseID   uuid.UUID  `json:"exercise_id"`
	WorkoutID    uuid.UUID  `json:"workout_id"`
	Weight       *int       `json:"weight"`
	Sets         *int       `json:"sets"`
	Reps         *int       `json:"reps"`
	Duration     *int       `json:"duration"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Completed    bool       `json:"completed"`
	ExerciseLoad int        `json:"exercise_load"`
}

func createdExerciseDTOFromDomain(
	exercise domain.CreatedExercise,
) CreatedExerciseDTOResponse {
	return CreatedExerciseDTOResponse{
		ID:           exercise.ID,
		Version:      exercise.Version,
		ExerciseID:   exercise.ExerciseID,
		WorkoutID:    exercise.WorkoutID,
		Weight:       exercise.Weight,
		Sets:         exercise.Sets,
		Reps:         exercise.Reps,
		Duration:     exercise.Duration,
		CreatedAt:    exercise.CreatedAt,
		UpdatedAt:    exercise.UpdatedAt,
		Completed:    exercise.Completed,
		ExerciseLoad: exercise.ExerciseLoad,
	}
}

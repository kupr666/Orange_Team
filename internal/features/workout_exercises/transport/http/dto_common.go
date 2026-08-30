package workout_exercises_transport_http

import (
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

type WorkoutExerciseDTOResponse struct {
	ID           uuid.UUID  `json:"id"`
	Version      int        `json:"version"`
	WorkoutID    uuid.UUID  `json:"workout_id"`
	ExerciseID   uuid.UUID  `json:"exercise_id"`
	Weight       *int       `json:"weight,omitempty"`
	Sets         *int       `json:"sets,omitempty"`
	Reps         *int       `json:"reps,omitempty"`
	Duration     *int       `json:"duration,omitempty"`
	Completed    bool       `json:"completed"`
	ExerciseLoad int        `json:"exercise_load"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func workoutExerciseDTOFromDomain(workoutExercise domain.WorkoutExercise) WorkoutExerciseDTOResponse {
	return WorkoutExerciseDTOResponse{
		ID:           workoutExercise.ID,
		Version:      workoutExercise.Version,
		WorkoutID:    workoutExercise.WorkoutID,
		ExerciseID:   workoutExercise.ExerciseID,
		Weight:       workoutExercise.Weight,
		Sets:         workoutExercise.Sets,
		Reps:         workoutExercise.Reps,
		Duration:     workoutExercise.Duration,
		Completed:    workoutExercise.Completed,
		ExerciseLoad: workoutExercise.ExerciseLoad,
		CreatedAt:    workoutExercise.CreatedAt,
		UpdatedAt:    workoutExercise.UpdatedAt,
	}
}

func workoutExerciseDTOsFromDomains(workoutExercises []domain.WorkoutExercise) []WorkoutExerciseDTOResponse {
	workoutExerciseDTO := make([]WorkoutExerciseDTOResponse, len(workoutExercises))
	for i, workoutExercise := range workoutExercises {
		workoutExerciseDTO[i] = workoutExerciseDTOFromDomain(workoutExercise)
	}
	return workoutExerciseDTO
}

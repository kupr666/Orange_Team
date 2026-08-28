package workout_exercises_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type WorkoutExercisesService interface {
	CreateWorkoutExercise(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
		exerciseID uuid.UUID,
		weight *int,
		sets *int,
		reps *int,
		duration *int,
		completed bool,
	) (domain.WorkoutExercise, error)

	GetWorkoutExercises(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) ([]domain.WorkoutExercise, error)

	PatchWorkoutExercise(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
		workoutExerciseID uuid.UUID,
		patch domain.WorkoutExercisePatch,
	) (domain.WorkoutExercise, error)

	DeleteWorkoutExercise(
		ctx context.Context,
		userID uuid.UUID,
		workoutExerciseID uuid.UUID,
		workoutID uuid.UUID,
	) error
}

type WorkoutExercisesHandler struct {
	workoutExercisesService WorkoutExercisesService
}

func NewWorkoutExercisesHandler(workoutExercisesService WorkoutExercisesService) *WorkoutExercisesHandler {
	return &WorkoutExercisesHandler{
		workoutExercisesService: workoutExercisesService,
	}
}

func (h *WorkoutExercisesHandler) Routes(
	authenticate core_http_middleware.Middleware,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/workouts/{workoutId}/exercises",
			Handler:    h.CreateWorkoutExercise,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodGet,
			Path:       "/workouts/{workoutId}/exercises",
			Handler:    h.GetWorkoutExercises,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/workouts/{workoutId}/exercises/{workoutExerciseId}",
			Handler:    h.PatchWorkoutExercise,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/workouts/{workoutId}/exercises/{workoutExerciseId}",
			Handler:    h.DeleteWorkoutExercise,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
	}
}

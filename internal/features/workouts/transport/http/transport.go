package workouts_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type WorkoutsService interface {
	// PostWorkout(
	// args
	// ) (retun values)

	GetWorkouts(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Workout, error)

	GetWorkout(
		ctx context.Context,
		workoutID uuid.UUID,
	) (domain.Workout, error)

	CreateExercise(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
		exerciseID uuid.UUID,
		weight *int,
		sets *int,
		reps *int,
		duration *int,
	) (domain.CreatedExercise, error)
	// PatchWorkout(
	// args
	// ) (retun values)

	// DeleteWorkout(
	// args
	// ) (retun values)
}

type WorkoutsHTTPHandler struct {
	workoutsService WorkoutsService
}

func NewWorkoutsHTTPHandler(workoutsService WorkoutsService) *WorkoutsHTTPHandler {
	return &WorkoutsHTTPHandler{
		workoutsService: workoutsService,
	}
}

func (h *WorkoutsHTTPHandler) Routes(
	authMiddleware core_http_middleware.Middleware,
) []core_http_server.Route {
	return []core_http_server.Route{
		// {
		// 	Method: http.MethodPost,
		// 	Path:   "/workouts",
		// 	Handler: h.CreateWorkout,
		// },
		{
			Method:  http.MethodGet,
			Path:    "/workouts",
			Handler: h.GetWorkouts,
		},
		{
			Method:  http.MethodGet,
			Path:    "/workouts/{workoutId}",
			Handler: h.GetWorkout,
		},
		{
			Method:     http.MethodPost,
			Path:       "/workouts/{workoutId}/exercises",
			Handler:    h.CreateExercise,
			Middleware: []core_http_middleware.Middleware{authMiddleware},
		},

		// {
		// 	Method: http.MethodPatch,
		// 	Path:   "/workouts/{workoutId}",
		// 	Handler: h.PatchWorkout,
		// },
		// {
		// 	Method: http.MethodDelete,
		// 	Path:   "/workouts/{workoutId}",
		// 	Handler: h.DeleteWorkout,
		// },
	}
}

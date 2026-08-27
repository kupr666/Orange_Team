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
	GetWorkouts(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Workout, error)

	GetWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) (domain.Workout, error)

	PatchWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
		patch domain.WorkoutPatch,
	) (domain.Workout, error)

	DeleteWorkout(
		ctx context.Context,
		userID uuid.UUID,
		workoutID uuid.UUID,
	) error
	CreateWorkout(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.Workout, error)

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
	authenticate core_http_middleware.Middleware,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/workouts",
			Handler: h.GetWorkouts,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodGet,
			Path:       "/workouts/{workoutId}",
			Handler:    h.GetWorkout,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/workouts/{workoutId}",
			Handler: h.PatchWorkout,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/workouts/{workoutId}",
			Handler: h.DeleteWorkout,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:  http.MethodPost,
			Path:    "/workouts",
			Handler: h.CreateWorkout,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
	}
}

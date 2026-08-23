package workouts_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
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
}

type WorkoutsHTTPHandler struct {
	workoutsService WorkoutsService
}

func NewWorkoutsHTTPHandler(workoutsService WorkoutsService) *WorkoutsHTTPHandler {
	return &WorkoutsHTTPHandler{
		workoutsService: workoutsService,
	}
}

func (h *WorkoutsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		// {
		// 	Method: http.MethodPost,
		// 	Path:   "/workouts",
		// 	// Handler: h.CreateWorkout,
		// },
		{
			Method: http.MethodGet,
			Path:   "/workouts",
			// Handler: h.GetWorkouts,
		},
		{
			Method: http.MethodGet,
			Path:   "/workouts/{workoutId}",
			// Handler: h.GetWorkout,
		},
		// {
		// 	Method: http.MethodPatch,
		// 	Path:   "/workouts/{workoutId}",
		// 	// Handler: h.PatchWorkout,
		// },
		// {
		// 	Method: http.MethodDelete,
		// 	Path:   "/workouts/{workoutId}",
		// 	// Handler: h.DeleteWorkout,
		// },
	}
}

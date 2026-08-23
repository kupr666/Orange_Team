package workouts_transport_http

import (
	"context"
	"net/http"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type WorkoutsService interface {
	GetWorkouts(
		ctx context.Context,
	//  
	) ([]domain.Workout, error)
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
		{
			Method: http.MethodGet,
			Path: "/workouts",
			// Handler: h.GetWorkouts,
		},
		{
			Method: http.MethodGet,
			Path: "/workouts/{id}",
			// Handler: h.GetWorkout,
		},
	}
}
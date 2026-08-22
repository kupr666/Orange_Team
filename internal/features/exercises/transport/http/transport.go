package exercises_transport_http

import (
	"context"
	"net/http"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type ExercisesService interface {
	GetExercises(
		ctx context.Context,
	) ([]domain.Exercise, error)
}

type ExercisesHTTPHandler struct {
	exercisesService ExercisesService
}

func NewExercisesHTTPHandler(exercisesService ExercisesService) *ExercisesHTTPHandler {
	return &ExercisesHTTPHandler{
		exercisesService: exercisesService,
	}
}

func (h *ExercisesHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/exercises",
			Handler: h.GetExercises,
		},
	}
}

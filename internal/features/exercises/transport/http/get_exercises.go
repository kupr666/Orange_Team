package exercises_transport_http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type ExerciseResponseDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Score       int        `json:"score"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type GetExercisesResponseDTO []ExerciseResponseDTO

func (h *ExercisesHTTPHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx, slog.Default())
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	exercises, err := h.exercisesService.GetExercises(ctx)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get exercises",
		)
		return
	}

	response := exercisesDTOFromDomain(exercises)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func exercisesDTOFromDomain(exercises []domain.Exercise) GetExercisesResponseDTO {
	response := make(GetExercisesResponseDTO, len(exercises))

	for i, exercise := range exercises {
		response[i] = ExerciseResponseDTO{
			ID:          exercise.ID,
			Name:        exercise.Name,
			Description: exercise.Description,
			Score:       exercise.Score,
			CreatedAt:   exercise.CreatedAt,
			UpdatedAt:   exercise.UpdatedAt,
		}
	}

	return response
}

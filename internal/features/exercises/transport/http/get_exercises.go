package exercises_transport_http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type ExerciseResponseDTO struct {
	ID          uuid.UUID        `json:"id"`
	Version     int64      `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  int        `json:"difficulty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Type        string     `json:"type"`
}

type GetExercisesResponseDTO []ExerciseResponseDTO

func (h *ExercisesHTTPHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
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

	responseHandler.JSONResponse(response, http.StatusOK)
}

func exercisesDTOFromDomain(exercises []domain.Exercise) GetExercisesResponseDTO {
	response := make(GetExercisesResponseDTO, len(exercises))

	for i, exercise := range exercises {
		response[i] = ExerciseResponseDTO{
			ID:          exercise.ID,
			Version:     exercise.Version,
			Name:        exercise.Name,
			Description: exercise.Description,
			Difficulty:  exercise.Difficulty,
			CreatedAt:   exercise.CreatedAt,
			UpdatedAt:   exercise.UpdatedAt,
			Type:        exercise.Type,
		}
	}

	return response
}

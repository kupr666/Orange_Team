package exercises_transport_http

import (
	"net/http"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateExerciseRequest struct {
	ID          int        `json:"id"`
	Version     int64      `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  int        `json:"difficulty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Type        string     `json:"type"`
}

type CreateExerciseResponse struct {
	ID          int        `json:"id"`
	Version     int64      `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  int        `json:"difficulty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Type        string     `json:"type"`
}

func (h *ExercisesHTTPHandler) CreateExercise(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateExerciseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	exerciseDomain := domain.NewExercise(
		request.ID,
		request.Version,
		request.Name,
		request.Description,
		request.Difficulty,
		request.CreatedAt,
		request.UpdatedAt,
		request.Type,
	)

	exerciseDomain, err := h.exercisesService.CreateExercise(ctx, exerciseDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create exercise",
		)

		return
	}

	response := exerciseDTOFromDomain(exerciseDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func exerciseDTOFromDomain(exercise domain.Exercise) CreateExerciseResponse{
	return CreateExerciseResponse{
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


package exercises_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateExerciseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  int    `json:"difficulty"`
	Type        string `json:"type"`
}

type CreateExerciseResponse ExerciseDTOResponse

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

	exerciseDomain, err := h.exercisesService.CreateExercise(
		ctx,
		request.Name,
		request.Description,
		request.Difficulty,
		request.Type,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create exercise",
		)

		return
	}

	response := CreateExerciseResponse(exerciseDTOFromDomain(exerciseDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

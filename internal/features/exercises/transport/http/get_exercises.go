package exercises_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type GetExercisesResponse []ExerciseDTOResponse

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

	response := GetExercisesResponse(exercisesDTOFromDomain(exercises))

	responseHandler.JSONResponse(response, http.StatusOK)
}

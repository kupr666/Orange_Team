package workouts_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *WorkoutsHTTPHandler) DeleteWorkout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workoutID path value",
		)

		return
	}

	if err := h.workoutsService.DeleteWorkout(ctx, workoutID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete workout",
		)

		return
	}

	responseHandler.NoContentResponse()
}

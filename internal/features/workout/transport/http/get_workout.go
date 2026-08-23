package workouts_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type GetWorkoutResponse WorkoutDTOResponse

func (h *WorkoutsHTTPHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workoutID path value",
		)

		return
	}

	workout, err := h.workoutsService.GetWorkout(ctx, workoutID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workout",
		)

		return
	}

	response := GetWorkoutResponse(workoutDTOFromDomain(workout))

	responseHandler.JSONResponse(response, http.StatusOK)
}

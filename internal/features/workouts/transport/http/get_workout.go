package workouts_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type GetWorkoutResponse WorkoutDTOResponse

func (h *WorkoutsHTTPHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	principal, ok := core_auth.PrincipalFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf(
				"authenticated principal is missing: %w",
				core_errors.ErrUnauthorized,
			),
			"valid JWT token is required",
		)
		return
	}

	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workoutID path value",
		)
		return
	}

	workout, err := h.workoutsService.GetWorkout(ctx, principal.UserID, workoutID)
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

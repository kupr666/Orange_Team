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

func (h *WorkoutsHTTPHandler) DeleteWorkout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

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

	if err := h.workoutsService.DeleteWorkout(ctx, principal.UserID, workoutID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete workout",
		)
		return
	}

	responseHandler.NoContentResponse()
}

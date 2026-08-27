package workouts_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *WorkoutsHTTPHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
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

	workouts, err := h.workoutsService.GetWorkouts(ctx, principal.UserID)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get workouts",
		)
		return
	}

	response := workoutDTOsFromDomains(workouts)

	responseHandler.JSONResponse(response, http.StatusOK)
}

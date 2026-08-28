package workout_exercises_transport_http

import (
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type GetWorkoutExercisesResponse []WorkoutExerciseDTOResponse

func (h *WorkoutExercisesHandler) GetWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	principal, ok := core_auth.PrincipalFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrUnauthorized,
			"authenticated user is missing",
		)
		return
	}

	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid workout ID",
		)
		return
	}

	workoutExerciseDomains, err := h.workoutExercisesService.GetWorkoutExercises(ctx, principal.UserID, workoutID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workout exercises",
		)
		return
	}

	response := GetWorkoutExercisesResponse(workoutExerciseDTOsFromDomains(workoutExerciseDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

// delete_workout_exercise.go (транспорт)
package workout_exercises_transport_http

import (
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *WorkoutExercisesHandler) DeleteWorkoutExercise(w http.ResponseWriter, r *http.Request) {
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

	workoutExerciseID, err := core_http_request.GetUUIDPathValue(r, "workoutExerciseId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid workout exercise ID",
		)
		return
	}

	if err := h.workoutExercisesService.DeleteWorkoutExercise(
		ctx,
		principal.UserID,
		workoutExerciseID,
		workoutID,
	); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete workout exercise",
		)
		return
	}

	responseHandler.NoContentResponse()
}

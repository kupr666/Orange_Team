package workout_exercises_transport_http

import (
	"fmt"
	"net/http"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
	core_http_types "github.com/kupr666/Orange_Team/internal/core/transport/http/types"
)

type PatchWorkoutExerciseRequest struct {
	Weight    core_http_types.Nullable[int]  `json:"weight"`
	Sets      core_http_types.Nullable[int]  `json:"sets"`
	Reps      core_http_types.Nullable[int]  `json:"reps"`
	Duration  core_http_types.Nullable[int]  `json:"duration"`
	Completed core_http_types.Nullable[bool] `json:"completed"`
}

func (r *PatchWorkoutExerciseRequest) Validate() error {
	if !r.Weight.Set && !r.Sets.Set && !r.Reps.Set && !r.Duration.Set && !r.Completed.Set {
		return fmt.Errorf("at least one field must be provided")
	}
	if r.Completed.Set && r.Completed.Value == nil {
		return fmt.Errorf("completed cannot be NULL")
	}
	return nil
}

type PatchWorkoutExerciseResponse WorkoutExerciseDTOResponse

func (h *WorkoutExercisesHandler) PatchWorkoutExercise(w http.ResponseWriter, r *http.Request) {
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
		responseHandler.ErrorResponse(err, "invalid workout ID")
		return
	}
	workoutExerciseID, err := core_http_request.GetUUIDPathValue(r, "workoutExerciseId")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid workout exercise ID")
		return
	}

	var request PatchWorkoutExerciseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	patch := domain.NewWorkoutExercisePatch(
		request.Weight.ToDomain(),
		request.Sets.ToDomain(),
		request.Reps.ToDomain(),
		request.Duration.ToDomain(),
		request.Completed.ToDomain(),
	)

	workoutExerciseDomain, err := h.workoutExercisesService.PatchWorkoutExercise(
		ctx,
		principal.UserID,
		workoutID,
		workoutExerciseID,
		patch,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch workout exercise",
		)
		return
	}

	response := PatchWorkoutExerciseResponse(workoutExerciseDTOFromDomain(workoutExerciseDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

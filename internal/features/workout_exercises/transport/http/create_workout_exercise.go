package workout_exercises_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateWorkoutExerciseRequest struct {
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
	Weight     *int      `json:"weight,omitempty"`
	Sets       *int      `json:"sets,omitempty"`
	Reps       *int      `json:"reps,omitempty"`
	Duration   *int      `json:"duration,omitempty"`
	Completed  bool      `json:"completed"`
}

type CreateWorkoutExerciseResponse WorkoutExerciseDTOResponse

func (h *WorkoutExercisesHandler) CreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
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

	var request CreateWorkoutExerciseRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	workoutExerciseDomain, err := h.workoutExercisesService.CreateWorkoutExercise(
		ctx,
		principal.UserID,
		workoutID,
		request.ExerciseID,
		request.Weight,
		request.Sets,
		request.Reps,
		request.Duration,
		request.Completed,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create workout exercise",
		)
		return
	}

	response := CreateWorkoutExerciseResponse(workoutExerciseDTOFromDomain(workoutExerciseDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

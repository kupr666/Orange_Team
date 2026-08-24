package workouts_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateExerciseRequest struct {
	ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`

	Weight *int `json:"weight,omitempty"`
	Sets   *int `json:"sets,omitempty"`
	Reps   *int `json:"reps,omitempty"`

	Duration *int `json:"duration,omitempty"` // секунды
}

func (h *WorkoutsHTTPHandler) CreateExercise(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// Пользователь, которого проверила JWT middleware.
	userID, ok := core_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.JSONResponse(
			core_http_response.ErrorResponse{
				Error:   "unauthorized",
				Message: "authenticated user is missing",
			},
			http.StatusUnauthorized,
		)
		return
	}

	// workoutId берём из URL, а не из JSON.
	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid workout ID")
		return
	}

	var request CreateExerciseRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "invalid create exercise request")
		return
	}

	createdExercise, err := h.workoutsService.CreateExercise(
		ctx,
		userID,
		workoutID,
		request.ExerciseID,
		request.Weight,
		request.Sets,
		request.Reps,
		request.Duration,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to add exercise to workout")
		return
	}

	response := createdExerciseDTOFromDomain(createdExercise)
	responseHandler.JSONResponse(response, http.StatusCreated)
}

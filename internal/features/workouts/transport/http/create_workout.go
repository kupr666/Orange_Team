package workouts_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateWorkoutResponse WorkoutDTOResponse

func (h *WorkoutsHTTPHandler) CreateWorkout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userID, ok := core_http_middleware.UserIDFromContext(ctx) // fix after merge branches
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
	createdWorkout, err := h.workoutsService.CreateWorkout(
		ctx,
		userID,
	)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create new workout")
		return
	}
	response := WorkoutDTOResponse(workoutDTOFromDomain(createdWorkout))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

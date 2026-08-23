package workouts_transport_http

import (
	"net/http"

	"github.com/google/uuid"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *WorkoutsHTTPHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	// get userID from context
	// path : request -> authentication middleware (token and session check) + 
	// put userID in r.Context - > we get this userID via helper function

	// userID, err := 

	// just because we haven't auth yet
	var userID uuid.UUID

	workouts, err := h.workoutsService.GetWorkouts(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to get workouts",
		)
		return
	}

	response := workoutDTOsFromDomains(workouts)

	responseHandler.JSONResponse(response, http.StatusOK)
}
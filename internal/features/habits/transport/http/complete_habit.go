package habits_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *HabitsHTTPHandler) CompleteHabit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), w)

	userID, err := authenticatedUserID(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "valid JWT token is required")
		return
	}
	habitID, err := core_http_request.GetUUIDPathValue(r, "habitId")
	if err != nil {
		responseHandler.ErrorResponse(err, "invalid habit ID")
		return
	}

	habit, err := h.habitsService.CompleteHabit(ctx, userID, habitID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to complete habit")
		return
	}

	responseHandler.JSONResponse(habitDTOFromDomain(habit), http.StatusOK)
}

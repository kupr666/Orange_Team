package habits_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *HabitsHTTPHandler) DeleteHabit(w http.ResponseWriter, r *http.Request) {
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

	if err := h.habitsService.DeleteHabit(ctx, userID, habitID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete habit")
		return
	}

	responseHandler.NoContentResponse()
}

package habits_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

func (h *HabitsHTTPHandler) GetHabits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), w)

	userID, err := authenticatedUserID(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "valid JWT token is required")
		return
	}

	habits, err := h.habitsService.GetHabits(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get habits")
		return
	}

	responseHandler.JSONResponse(habitsDTOFromDomain(habits), http.StatusOK)
}

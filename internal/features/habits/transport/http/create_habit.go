package habits_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateHabitRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=80"`
	Description string `json:"description" validate:"max=500"`
}

type CreateHabitResponse HabitDTOResponse

func (h *HabitsHTTPHandler) CreateHabit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := core_http_response.NewHTTPResponseHandler(core_logger.FromContext(ctx), w)

	userID, err := authenticatedUserID(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "valid JWT token is required")
		return
	}

	var request CreateHabitRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate habit")
		return
	}

	habit, err := h.habitsService.CreateHabit(ctx, userID, request.Name, request.Description)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create habit")
		return
	}

	responseHandler.JSONResponse(CreateHabitResponse(habitDTOFromDomain(habit)), http.StatusCreated)
}

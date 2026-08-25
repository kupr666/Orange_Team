package authentication_transport_http

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type RegisterUserRequestDTO struct {
	Email    string `json:"email" validate:"required,email,max=30"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2,max=50"`
}

type RegisterUserResponseDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *AuthenticationHTTPHandler) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request RegisterUserRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	createdUser, err := h.authenticationService.RegisterUser(
		ctx,
		request.Email,
		request.Password,
		request.FullName,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create user",
		)
		return
	}

	response := userDTOFromDomain(createdUser)

	responseHandler.JSONResponse(response, http.StatusCreated)
}

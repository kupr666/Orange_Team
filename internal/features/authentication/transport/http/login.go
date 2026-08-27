package authentication_transport_http

import (
	"net/http"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type loginRequestDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponseDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *AuthenticationHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request loginRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	loginResult, err := h.authenticationService.Login(
		ctx,
		request.Email,
		request.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to login",
		)

		return
	}

	response := loginResponseDTO{
		AccessToken: loginResult.AccessToken,
		TokenType:   loginResult.TokenType,
	}
	responseHandler.JSONResponse(response, http.StatusOK)
}

package authentication_transport_http

import (
	"context"
	"net/http"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type AuthenticationService interface {
	RegisterUser(
		ctx context.Context,
		email string,
		password string,
		fullName string,
	) (domain.User, error)
}

type AuthenticationHTTPHandler struct {
	authenticationService AuthenticationService
}

func NewAuthenticationHTTPHandler(
	authenticationService AuthenticationService,
) *AuthenticationHTTPHandler {
	return &AuthenticationHTTPHandler{
		authenticationService: authenticationService,
	}
}

func (h *AuthenticationHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.RegisterUser,
		},
	}
}

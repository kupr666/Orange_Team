package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type AuthService interface {
	RegisterUser(
		ctx context.Context,
		email string,
		password string,
		fullName string,
	) (domain.User, error)

	LoginUser(
		ctx context.Context,
		email string,
		password string,
	) (domain.User, error)
}

type AuthHTTPHandler struct {
	authService AuthService
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.RegisterUser,
		},
		// {
		// 	Method:  http.MethodPost,
		// 	Path:    "/login",
		// 	Handler: h.LoginUser,
		// },
	}
}

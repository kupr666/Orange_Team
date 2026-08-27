package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type UsersService interface {
	GetUser(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.User, error)

	PatchUser(
		ctx context.Context,
		userID uuid.UUID,
		patch domain.UserPatch,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type UsersHTTPHandler struct {
	usersService UsersService
}

func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes(
	authentication core_http_middleware.Middleware,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/users/me",
			Handler:    h.GetUser,
			Middleware: []core_http_middleware.Middleware{authentication},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/users/me",
			Handler:    h.PatchUser,
			Middleware: []core_http_middleware.Middleware{authentication},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/users/me",
			Handler:    h.DeleteUser,
			Middleware: []core_http_middleware.Middleware{authentication},
		},
	}
}

package users_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
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

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/me",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/me",
			Handler: h.PatchUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/me",
			Handler: h.DeleteUser,
		},
	}
}

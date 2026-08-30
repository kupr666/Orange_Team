package habits_transport_http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_middleware "github.com/kupr666/Orange_Team/internal/core/transport/http/middleware"
	core_http_server "github.com/kupr666/Orange_Team/internal/core/transport/http/server"
)

type HabitsService interface {
	CreateHabit(ctx context.Context, userID uuid.UUID, name string, description string) (domain.Habit, error)
	GetHabits(ctx context.Context, userID uuid.UUID) ([]domain.Habit, error)
	CompleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) (domain.Habit, error)
	DeleteHabit(ctx context.Context, userID uuid.UUID, habitID uuid.UUID) error
}

type HabitsHTTPHandler struct {
	habitsService HabitsService
}

func NewHabitsHTTPHandler(habitsService HabitsService) *HabitsHTTPHandler {
	return &HabitsHTTPHandler{habitsService: habitsService}
}

func (h *HabitsHTTPHandler) Routes(
	authenticate core_http_middleware.Middleware,
) []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/habits",
			Handler:    h.GetHabits,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodPost,
			Path:       "/habits",
			Handler:    h.CreateHabit,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodPost,
			Path:       "/habits/{habitId}/completions",
			Handler:    h.CompleteHabit,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/habits/{habitId}",
			Handler:    h.DeleteHabit,
			Middleware: []core_http_middleware.Middleware{authenticate},
		},
	}
}

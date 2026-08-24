package workouts_transport_http

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
)

type CreateWorkoutRequestDTO struct {
	PersonalScoreCoefficient int `json:"personal_score_coefficient"`
}

type WorkoutResponseDTO struct {
	ID                       uuid.UUID  `json:"id"`
	Status                   string     `json:"status"`
	StartedAt                *time.Time `json:"started_at"`
	CompletedAt              *time.Time `json:"completed_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	WorkoutScore             int        `json:"workout_score"`
	Intensity                *int       `json:"intensity"`
	PersonalScoreCoefficient int        `json:"personal_score_coefficient"`
}

func (h *WorkoutsHTTPHandler) CreateWorkout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateWorkoutRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	
}

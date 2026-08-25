package workouts_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
	core_http_types "github.com/kupr666/Orange_Team/internal/core/transport/http/types"
)

type PatchWorkoutRequest struct {
	Status      core_http_types.Nullable[string] `json:"status"`
	StartedAt   core_http_types.Nullable[string] `json:"started_at"`
	CompletedAt core_http_types.Nullable[string] `json:"completed_at"`
	Intensity   core_http_types.Nullable[int]    `json:"intensity"`
}

func (r *PatchWorkoutRequest) Validate() error {
	if !r.Status.Set && !r.StartedAt.Set && !r.CompletedAt.Set && !r.Intensity.Set {
		return fmt.Errorf("at least one field must be provided")
	}

	if r.Status.Set {
		if r.Status.Value == nil {
			return fmt.Errorf("`Status` cannot be NULL")
		}
		status := *r.Status.Value

		if !domain.AllowedStatuses[status] {
			return fmt.Errorf("`Status` must be one of: planned, in_progress, completed, cancelled")
		}
	}

	if r.StartedAt.Set && r.StartedAt.Value != nil {
		if _, err := time.Parse(time.RFC3339, *r.StartedAt.Value); err != nil {
			return fmt.Errorf("`StartedAt` must be a string in RFC3339 format")
		}
	}
	if r.CompletedAt.Set && r.CompletedAt.Value != nil {
		if _, err := time.Parse(time.RFC3339, *r.CompletedAt.Value); err != nil {
			return fmt.Errorf("`CompletedAt` must be a string in RFC3339 format")
		}
	}

	if r.Intensity.Set && r.Intensity.Value != nil {
		intensity := *r.Intensity.Value
		if intensity < 1 || intensity > 10 {
			return fmt.Errorf("`Intensity` must be between 1 and 10 or NULL")
		}
	}

	return nil
}

type PatchWorkoutResponse WorkoutDTOResponse

func (h *WorkoutsHTTPHandler) PatchWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	workoutID, err := core_http_request.GetUUIDPathValue(r, "workoutId")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get workout ID from path",
		)
		return
	}

	var request PatchWorkoutRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	workoutPatch, err := workoutPatchFromRequest(request)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid request data",
		)
		return
	}

	updatedWorkout, err := h.workoutsService.PatchWorkout(ctx, workoutID, workoutPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch workout",
		)
		return
	}

	response := PatchWorkoutResponse(workoutDTOFromDomain(updatedWorkout))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func workoutPatchFromRequest(request PatchWorkoutRequest) (domain.WorkoutPatch, error) {
	startedAt, err := toDomainNullableTime(request.StartedAt)
	if err != nil {
		return domain.WorkoutPatch{}, fmt.Errorf("parse started_at: %w", err)
	}
	completedAt, err := toDomainNullableTime(request.CompletedAt)
	if err != nil {
		return domain.WorkoutPatch{}, fmt.Errorf("parse completed_at: %w", err)
	}

	return domain.NewWorkoutPatch(
		request.Status.ToDomain(),
		startedAt,
		completedAt,
		request.Intensity.ToDomain(),
	), nil
}

func toDomainNullableTime(nullable core_http_types.Nullable[string]) (domain.Nullable[time.Time], error) {
	if !nullable.Set {
		return domain.Nullable[time.Time]{Set: false}, nil
	}
	if nullable.Value == nil {
		return domain.Nullable[time.Time]{Value: nil, Set: true}, nil
	}
	t, err := time.Parse(time.RFC3339, *nullable.Value)
	if err != nil {
		return domain.Nullable[time.Time]{}, fmt.Errorf("invalid RFC3339 time: %w", err)
	}
	return domain.Nullable[time.Time]{Value: &t, Set: true}, nil
}

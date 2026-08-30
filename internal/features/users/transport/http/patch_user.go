package users_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_auth "github.com/kupr666/Orange_Team/internal/core/auth"
	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_errors "github.com/kupr666/Orange_Team/internal/core/errors"
	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
	core_http_request "github.com/kupr666/Orange_Team/internal/core/transport/http/request"
	core_http_response "github.com/kupr666/Orange_Team/internal/core/transport/http/response"
	core_http_types "github.com/kupr666/Orange_Team/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	Sex         core_http_types.Nullable[string] `json:"sex"`
	WeightGrams core_http_types.Nullable[int]    `json:"weight_grams"`
	BirthDate   core_http_types.Nullable[string] `json:"birth_date"`
	HeightCM    core_http_types.Nullable[int]    `json:"height_cm"`
}

func (r *PatchUserRequest) Validate() error {
	if !r.Sex.Set && !r.WeightGrams.Set && !r.BirthDate.Set && !r.HeightCM.Set {
		return fmt.Errorf("at least one field must be provided")
	}

	if r.Sex.Set && r.Sex.Value != nil {
		if !domain.AllowedSexes[*r.Sex.Value] {
			return fmt.Errorf(
				"sex must be one of: %s, %s",
				domain.SexMale,
				domain.SexFemale)
		}
	}

	if r.WeightGrams.Set && r.WeightGrams.Value != nil {
		if *r.WeightGrams.Value < domain.MinWeightGrams || *r.WeightGrams.Value > domain.MaxWeightGrams {
			return fmt.Errorf("weight_grams must be between %d and %d",
				domain.MinWeightGrams,
				domain.MaxWeightGrams,
			)
		}
	}

	if r.BirthDate.Set && r.BirthDate.Value != nil {
		if _, err := time.Parse("2006-01-02", *r.BirthDate.Value); err != nil {
			return fmt.Errorf("birth_date must be in YYYY-MM-DD format")
		}
	}

	if r.HeightCM.Set && r.HeightCM.Value != nil {
		if *r.HeightCM.Value < domain.MinHeightCM || *r.HeightCM.Value > domain.MaxHeightCM {
			return fmt.Errorf("height_cm must be between %d and %d",
				domain.MinHeightCM,
				domain.MaxHeightCM,
			)
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	principal, ok := core_auth.PrincipalFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			core_errors.ErrUnauthorized,
			"authenticated user is missing",
		)
		return
	}
	userID := principal.UserID

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	patch, err := userPatchFromRequest(request)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"invalid request data",
		)
		return
	}

	updatedUser, err := h.usersService.PatchUser(ctx, userID, patch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(updatedUser))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) (domain.UserPatch, error) {
	birthDate, err := core_http_request.ToDomainNullableDate(request.BirthDate)
	if err != nil {
		return domain.UserPatch{}, fmt.Errorf(
			"parse birth_date: %w",
			err,
		)
	}

	return domain.NewUserPatch(
		request.Sex.ToDomain(),
		request.WeightGrams.ToDomain(),
		birthDate,
		request.HeightCM.ToDomain(),
	), nil
}

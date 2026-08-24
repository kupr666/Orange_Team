package authentication_transport_http

import "github.com/kupr666/Orange_Team/internal/core/domain"

func userDTOFromDomain(user domain.User) RegisterUserResponseDTO {
	return RegisterUserResponseDTO{
		ID:               user.ID,
		Email:            user.Email,
		FullName:         user.FullName,
		ProfileCompleted: user.ProfileCompleted(),
		CreatedAt:        user.CreatedAt,
	}
}

package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kupr666/Orange_Team/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID uuid.UUID,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"get user from repository: %w",
			err,
		)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf(
			"apply user patch: %w",
			err,
		)
	}

	updatedUser, err := s.usersRepository.PatchUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"update user in repository: %w",
			err,
		)
	}

	return updatedUser, nil
}

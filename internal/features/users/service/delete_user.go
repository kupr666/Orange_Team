package users_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) error {
	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return  nil
}

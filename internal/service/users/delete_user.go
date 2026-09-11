package users

import (
	"context"
	"fmt"
)

func (s *service) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

package users

import (
	"context"
	"fmt"
)

func (s *service) UpdateUser(ctx context.Context, userID, login string) error {
	err := s.userRepo.UpdateUser(ctx, userID, login)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

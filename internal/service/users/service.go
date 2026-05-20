package users

import (
	"context"
	"fmt"
)

type service struct {
	userRepo userRepository
}

func New(userRepo userRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) GetUser(ctx context.Context, userID string) (string, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, userID, username string) error {
	err := s.userRepo.UpdateUser(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *service) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

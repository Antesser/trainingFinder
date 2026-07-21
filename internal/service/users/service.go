package users

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/user"
)

type service struct {
	userRepo userRepository
}

func New(userRepo userRepository) *service {
	return &service{userRepo: userRepo}
}

func (s *service) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, userID, login string) error {
	err := s.userRepo.UpdateUser(ctx, userID, login)
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

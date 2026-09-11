package users

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/user"
)

func (s *service) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(ctx context.Context, login, password, roleID string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	userId := uuid.New().String()
	id, err := s.authRepo.SignUp(ctx, hash, userId, login, roleID)
	if err != nil {
		return "", fmt.Errorf("sign up user: %w", err)
	}
	return id, nil
}

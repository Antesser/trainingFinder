package auth

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(ctx context.Context, login, password string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userId := uuid.New().String()
	id, err := s.authRepo.SignUp(ctx, hash, userId, login)
	if err != nil {
		return "", err
	}
	return id, nil
}

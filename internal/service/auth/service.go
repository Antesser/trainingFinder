package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type service struct {
	authRepo authRepository
}

func New(authRepo authRepository) *service {
	return &service{
		authRepo: authRepo,
	}
}
func generateRandomHash(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *service) SignUp(ctx context.Context, login, password string) (string, string, error) {
	hash, _ := generateRandomHash(16)
	log, pass, err := s.authRepo.SignUp(ctx, hash, login, password)
	if err != nil {
		return "", "", err
	}
	return log, pass, nil
}

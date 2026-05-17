package auth

import "context"

type service struct {
	authRepo authRepository
}

func New(authRepo authRepository) *service {
	return &service{
		authRepo: authRepo,
	}
}

func (s *service) SignUp(ctx context.Context, login, password string) (string, string, error) {
	log, pass, err := s.authRepo.SignUp(ctx, login, password)
	if err != nil {
		return "", "", err
	}
	return log, pass, nil
}

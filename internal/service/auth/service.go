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
	// bussines logic
	// s.authRepo.
	// call repostory (interface)

	return "", "", nil
}

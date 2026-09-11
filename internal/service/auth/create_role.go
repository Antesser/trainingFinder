package auth

import "context"

func (s *service) CreateRole(ctx context.Context, role string) error {
	err := s.authRepo.CreateRole(ctx, role)
	if err != nil {
		return err
	}
	return nil
}

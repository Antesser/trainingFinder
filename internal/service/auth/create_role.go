package auth

import (
	"context"
	"fmt"
)

func (s *service) CreateRole(ctx context.Context, role string) error {
	if err := s.authRepo.CreateRole(ctx, role); err != nil {
		return fmt.Errorf("create role %q: %w", role, err)
	}
	return nil
}

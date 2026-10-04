package users

import (
	"context"

	userPkg "github.com/Antesser/trainingFinder/pkg/api/users/v1"
)

func (s *Server) DeleteUser(ctx context.Context, req *userPkg.DeleteUserRequest) (*userPkg.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userPkg.DeleteUserResponse{}, nil
}

package auth

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authService "trainingFinder/internal/service/auth"
	authPkg "trainingFinder/pkg/api/auth"
)

type Server struct {
	authPkg.UnimplementedAuthServiceServer
	authService AuthService
}

func (s *Server) SignUp(ctx context.Context, req *authPkg.SignUpRequest) (*authPkg.SignUpResponse, error) {
	log.Printf("SignUp request: login=%s", req.Login)

	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password required")
	}

	// call service (interface)

	accessToken := "some-access-token"
	refreshToken := "some-refresh-token"

	return &authPkg.SignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

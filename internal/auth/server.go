package auth

import (
    "context"
    "log"

    authv1 "trainingFinder/pkg/api/auth" 
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type Server struct {
    authv1.UnimplementedAuthServiceServer
}

func (s *Server) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
    log.Printf("SignUp request: login=%s", req.Login)

    if req.Login == "" || req.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "login and password required")
    }

    accessToken := "some-access-token"
    refreshToken := "some-refresh-token"

    return &authv1.SignUpResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
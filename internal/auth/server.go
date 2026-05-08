package auth

import (
    "context"
    "log"

    authPkg "trainingFinder/pkg/api/auth"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type Server struct {
    authPkg.UnimplementedAuthServiceServer
}

func (s *Server) SignUp(ctx context.Context, req *authPkg.SignUpRequest) (*authPkg.SignUpResponse, error) {
    log.Printf("SignUp request: login=%s", req.Login)

    if req.Login == "" || req.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "login and password required")
    }

    accessToken := "some-access-token"
    refreshToken := "some-refresh-token"

    return &authPkg.SignUpResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
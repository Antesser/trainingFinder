package auth

import (
	"context"
	"log"

	authPkg "trainingFinder/pkg/api/auth/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct { // структура для домена auth и тд
	authPkg.UnimplementedAuthServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	AuthService                            authService
}

func NewServer(s authService) *Server {
	return &Server{AuthService: s}
}

func (s *Server) RegisterServer(server *grpc.Server) {
	authPkg.RegisterAuthServiceServer(server, s)
}

func (s *Server) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	err := authPkg.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
	if err != nil {
		log.Fatal("Registration error", err)
	}

	return err
}

func (s *Server) SignUp(ctx context.Context, req *authPkg.SignUpRequest) (*authPkg.SignUpResponse, error) {
	log.Printf("SignUp request: login=%s", req.Login)

	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password required")
	}

	accessToken, refreshToken, err := s.AuthService.SignUp(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &authPkg.SignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

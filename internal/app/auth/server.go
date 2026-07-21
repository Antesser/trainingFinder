package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	authPkg "github.com/Antesser/trainingFinder/pkg/api/auth/v1"

	model "github.com/Antesser/trainingFinder/internal/model/training"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct { // структура для домена auth и тд
	authPkg.UnimplementedAuthServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	authService                            authService
}

func NewServer(s authService) *Server {
	return &Server{authService: s}
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

func (s *Server) SignIn(ctx context.Context, req *authPkg.SignInRequest) (*authPkg.SignInResponse, error) {
	tokens, err := s.authService.SignIn(ctx, req.Login, req.Password)
	if err != nil {

		if errors.Is(err, model.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, err
	}
	// добавить печеньки, в которые я положу refreshToken, проблема в том, что всё может пойти по ...
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	cookieStr := cookie.String()
	header := metadata.Pairs("Set-Cookie", cookieStr)
	grpc.SendHeader(ctx, header)
	return &authPkg.SignInResponse{
		AccessToken: tokens.AccessToken,
	}, nil
}
func (s *Server) SignUp(ctx context.Context, req *authPkg.SignUpRequest) (*authPkg.SignUpResponse, error) { // вынести в отдельный файл, Виталий негодует
	log.Printf("SignUp request: login=%s", req.Login)

	id, err := s.authService.SignUp(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &authPkg.SignUpResponse{
		Id: id,
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *authPkg.RefreshTokenRequest) (*authPkg.RefreshTokenResponse, error) {
	refToken, err := uuid.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.authService.RefreshSession(ctx, refToken)
	if err != nil {
		return nil, err
	}

	return &authPkg.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}
func (s *Server) CreateRole(ctx context.Context, req *authPkg.CreateRoleRequest) (*authPkg.CreateRoleResponse, error) {

	err := s.authService.CreateRole(ctx, req.Role)
	if err != nil {
		return nil, err
	}

	return &authPkg.CreateRoleResponse{}, nil
}

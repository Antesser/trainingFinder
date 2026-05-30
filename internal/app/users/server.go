package users

import (
	"context"
	"log"
	userPkg "trainingFinder/pkg/api/users/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct { // структура для домена auth и тд
	userPkg.UnimplementedUserServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	UserService                            userService
}

func NewServer(s userService) *Server {
	return &Server{UserService: s}
}

func (s *Server) RegisterServer(server *grpc.Server) {
	userPkg.RegisterUserServiceServer(server, s)
}

func (s *Server) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	err := userPkg.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
	if err != nil {
		log.Fatal("Registration error", err)
	}

	return err
}

func (s *Server) DeleteUser(ctx context.Context, req *userPkg.DeleteUserRequest) (*userPkg.DeleteUserResponse, error) {

	err := s.UserService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userPkg.DeleteUserResponse{}, nil
}
func (s *Server) UpdateUser(ctx context.Context, req *userPkg.UpdateUserRequest) (*userPkg.UpdateUserResponse, error) { //вынести в отдельные файлы, ибо надо

	err := s.UserService.UpdateUser(ctx, req.GetId(), req.GetUserByIDname())
	if err != nil {
		return nil, err
	}

	return &userPkg.UpdateUserResponse{}, nil
}
func (s *Server) GetUserByID(ctx context.Context, req *userPkg.GetUserByIDRequest) (*userPkg.GetUserByIDResponse, error) {

	//id,login, err := s.UserService.GetUserByID(ctx, req.GetId())
	//if err != nil {
	//	return nil, err
	//}

	return &userPkg.GetUserByIDResponse{Id: "sdasg", Login: "sdkjgfh"}, nil
}

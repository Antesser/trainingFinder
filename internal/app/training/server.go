package training

import (
	"context"
	"log"

	trainingPkg "trainingFinder/pkg/api/training/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct { // структура для домена auth и тд
	trainingPkg.UnimplementedTrainingServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	AuthService                                    trainingService
}

func NewServer(s trainingService) *Server {
	return &Server{AuthService: s}
}

func (s *Server) RegisterServer(server *grpc.Server) {
	trainingPkg.RegisterTrainingServiceServer(server, s)
}

func (s *Server) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	err := trainingPkg.RegisterTrainingServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
	if err != nil {
		log.Fatal("Registration error", err)
	}

	return err
}

func (s *Server) CreateTraining(ctx context.Context, req *trainingPkg.CreateTrainingRequest) (*trainingPkg.CreateTrainingResponse, error) {
	//string trainer_id = 1 [(buf.validate.field).string.uuid = true];
	//string user_id = 2 [(buf.validate.field).required = true];
	//google.protobuf.Timestamp started_at = 3; // указываем полный путь до структуры
	//google.protobuf.Timestamp ended_at = 4;
	//string additional_info = 5;
	accessToken, refreshToken, err := s.AuthService.SignUp(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	return &trainingPkg.CreateTrainingResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

package training

import (
	"context"
	"log"
	model "trainingFinder/internal/model/training"

	trainingPkg "trainingFinder/pkg/api/training/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct { // структура для домена
	trainingPkg.UnimplementedTrainingServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	TrainingService                                trainingService
}

func NewServer(s trainingService) *Server {
	return &Server{TrainingService: s}
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

	trainingModel := model.Training{
		TrainerID:      req.TrainerId,
		UserID:         req.UserId,
		StartedAt:      req.StartedAt.AsTime(),
		EndedAt:        req.EndedAt.AsTime(),
		AdditionalInfo: req.AdditionalInfo,
	}
	err := s.TrainingService.CreateTraining(ctx, &trainingModel)
	if err != nil {
		return nil, err
	}

	return &trainingPkg.CreateTrainingResponse{}, nil
}

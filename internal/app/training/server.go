package training

import (
	"context"
	"errors"
	"log"
	model "trainingFinder/internal/model/training"

	trainingPkg "trainingFinder/pkg/api/training/v1"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct { // структура для домена
	trainingPkg.UnimplementedTrainingServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	trainingService                                trainingService
}

func NewServer(s trainingService) *Server {
	return &Server{trainingService: s}
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
		ID:             uuid.New().String(),
		TrainerID:      req.TrainerId,
		UserID:         req.UserId,
		StartedAt:      req.StartedAt.AsTime(),
		EndedAt:        req.EndedAt.AsTime(),
		AdditionalInfo: req.AdditionalInfo,
	}
	err := s.trainingService.CreateTraining(ctx, &trainingModel)
	if err != nil {
		return nil, err
	}

	return &trainingPkg.CreateTrainingResponse{}, nil
}
func (s *Server) GetTraining(ctx context.Context, req *trainingPkg.GetTrainingRequest) (*trainingPkg.GetTrainingResponse, error) {

	mod, err := s.trainingService.GetTraining(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &trainingPkg.GetTrainingResponse{Training: toProtoTraining(*mod)},
		nil
}
func (s *Server) UpdateTraining(ctx context.Context, req *trainingPkg.UpdateTrainingRequest) (*trainingPkg.UpdateTrainingResponse, error) {
	trainingModel := model.UpdateTrainingRequest{
		ID:        req.Id,
		TrainerID: req.TrainerId,
		UserID:    req.UserId,
		//StartedAt:      req.StartedAt.AsTime(),
		//EndedAt:        req.EndedAt.AsTime(),
		AdditionalInfo: req.AdditionalInfo,
	}
	if req.StartedAt != nil {
		trainingModel.StartedAt = new(req.StartedAt.AsTime())
	}
	if req.EndedAt != nil {
		trainingModel.EndedAt = new(req.EndedAt.AsTime())
	}
	err := s.trainingService.UpdateTraining(ctx, trainingModel)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &trainingPkg.UpdateTrainingResponse{},
		nil
}
func (s *Server) DeleteTraining(ctx context.Context, req *trainingPkg.DeleteTrainingRequest) (*trainingPkg.DeleteTrainingResponse, error) {
	err := s.trainingService.DeleteTraining(ctx, req.Id)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &trainingPkg.DeleteTrainingResponse{},
		nil
}
func (s *Server) ListTrainings(ctx context.Context, _ *trainingPkg.ListTrainingsRequest) (*trainingPkg.ListTrainingsResponse, error) {

	mod, err := s.trainingService.ListTraining(ctx)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	protoList := make([]*trainingPkg.Training, len(mod))
	for i, model := range mod {
		protoList[i] = toProtoTraining(model)
	}

	return &trainingPkg.ListTrainingsResponse{
		Trainings: protoList,
	}, nil
}

func (s *Server) BookTraining(ctx context.Context, req *trainingPkg.BookTrainingRequest) (*trainingPkg.BookTrainingResponse, error) {
	err := s.trainingService.BookTraining(ctx, req.TrainingID, req.UserID)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &trainingPkg.BookTrainingResponse{},
		nil
}

func toProtoTraining(m model.Training) *trainingPkg.Training {
	return &trainingPkg.Training{
		Id:             m.ID,
		TrainerId:      m.TrainerID,
		UserId:         m.UserID,
		StartedAt:      timestamppb.New(m.StartedAt),
		EndedAt:        timestamppb.New(m.EndedAt),
		AdditionalInfo: m.AdditionalInfo,
	}
}

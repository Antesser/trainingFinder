package booking

import (
	"context"
	"errors"
	"log"
	model "trainingFinder/internal/model/training"
	bookingPkg "trainingFinder/pkg/api/booking/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct { // структура для домена
	bookingPkg.UnimplementedBookingServiceServer // реализуем имплементацию, хранящуюся внутри grpc (которая там сгенерирована, не реализована)
	bookingService                               bookingService
}

func NewServer(s bookingService) *Server {
	return &Server{bookingService: s}
}

func (s *Server) RegisterServer(server *grpc.Server) {
	bookingPkg.RegisterBookingServiceServer(server, s)
}

func (s *Server) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	err := bookingPkg.RegisterBookingServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
	if err != nil {
		log.Fatal("Registration error", err)
	}

	return err
}

// todo вынести в отдельный модуль или сделать сабмодулем training u know?
func (s *Server) BookTraining(ctx context.Context, req *bookingPkg.BookingTrainingRequest) (*bookingPkg.BookingTrainingResponse, error) {
	err := s.bookingService.BookTraining(ctx, req.TrainingId, req.UserId)
	if err != nil {
		if errors.Is(err, model.ErrTrainingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &bookingPkg.BookingTrainingResponse{},
		nil
}

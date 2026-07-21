package booking

import (
	"context"
	"errors"
	"log"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	modelPage "github.com/Antesser/trainingFinder/internal/model/page"
	bookingPkg "github.com/Antesser/trainingFinder/pkg/api/booking/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func (s *Server) BookTraining(ctx context.Context, req *bookingPkg.BookingTrainingRequest) (*bookingPkg.BookingTrainingResponse, error) {
	err := s.bookingService.BookTraining(ctx, model.TrainingBooking{TrainingID: req.TrainingId, UserID: req.UserId, BookFrom: req.BookFrom.AsTime(), BookTo: req.BookTo.AsTime()})
	if err != nil {
		if errors.Is(err, model.ErrBookingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &bookingPkg.BookingTrainingResponse{},
		nil
}
func (s *Server) ListTraining(ctx context.Context, req *bookingPkg.ListBookingsRequest) (*bookingPkg.ListBookingsResponse, error) {
	mod, hasNext, err := s.bookingService.ListBookings(ctx, modelPage.Page{Limit: req.Page.Limit, Offset: req.Page.Offset}, req.BookedBy)
	if err != nil {
		if errors.Is(err, model.ErrBookingNotFound) {
			return &bookingPkg.ListBookingsResponse{
				Bookings: []*bookingPkg.Booking{},
				HasNext:  hasNext,
			}, nil
		}
		return nil, err
	}

	protoList := lo.Map(mod, func(m model.TrainingBooking, _ int) *bookingPkg.Booking {
		return toProtoTraining(&m)
	})

	return &bookingPkg.ListBookingsResponse{
		Bookings: protoList,
		HasNext:  hasNext,
	}, nil
}

func toProtoTraining(m *model.TrainingBooking) *bookingPkg.Booking {
	return &bookingPkg.Booking{
		BookedBy:   m.UserID,
		TrainingId: m.TrainingID,
		BookFrom:   timestamppb.New(m.BookFrom),
		BookTo:     timestamppb.New(m.BookTo),
	}
}

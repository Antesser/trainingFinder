package booking

import (
	"context"

	modelBooking "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
	modelTraining "github.com/Antesser/trainingFinder/internal/model/training"
)

type bookingRepository interface {
	GetStatus(ctx context.Context, bookingID string) (modelBooking.Status, error)
	ChangeStatus(ctx context.Context, data modelBooking.BookingStatus) error
	CreateTrainingBooking(ctx context.Context, mod modelBooking.TrainingBooking) (string, error)
	ListBookings(ctx context.Context, data modelBooking.ListBookingsRequest) (*modelBooking.ListBookingsResponse, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
}

type trainingBookingRepository interface {
	GetTrainingByBookingID(ctx context.Context, bookingID string) (*modelTraining.Training, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

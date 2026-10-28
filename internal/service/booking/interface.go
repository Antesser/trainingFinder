package booking

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/booking"
	modelBooking "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
	modelTraining "github.com/Antesser/trainingFinder/internal/model/training"
	"github.com/jackc/pgx/v5"
)

type bookingRepository interface {
	GetStatus(ctx context.Context, bookingID string) (modelBooking.Status, error)
	ChangeStatus(ctx context.Context, data modelBooking.BookingStatus) error
	CreateTrainingBooking(ctx context.Context, mod modelBooking.TrainingBooking) error
	ListBookings(ctx context.Context, data booking.ListBookingsRequest) (*modelBooking.ListBookingsResponse, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
}
type trainingRepository interface {
	GetTraining(ctx context.Context, id string) (*modelTraining.Training, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	InTxWithIsoLevel(ctx context.Context, isoLevel pgx.TxIsoLevel, fn func(ctx context.Context) error) error
}

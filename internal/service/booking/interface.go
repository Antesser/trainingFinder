package booking

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/booking"
	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
)

type bookingRepository interface {
	CreateTrainingBooking(ctx context.Context, mod model.TrainingBooking, withLock bool) error
	CheckIntersections(ctx context.Context, booking model.TrainingBooking) error
	ListBookings(ctx context.Context, data booking.ListBookingsRequest) (model.ListBookingsResponse, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
}

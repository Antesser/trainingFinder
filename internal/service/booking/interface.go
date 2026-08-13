package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

type bookingRepository interface {
	CreateTrainingBooking(ctx context.Context, mod model.TrainingBooking) error
	ListBookings(ctx context.Context, page page.Page, bookedBy string) ([]model.TrainingBooking, bool, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
}

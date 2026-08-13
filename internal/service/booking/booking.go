package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

type service struct {
	repo              bookingRepository
	outboxRepo        outboxRepository
	bookingMarshaller bookingMarshaller
}
type bookingMarshaller func(event model.CreateBookingEvent) ([]byte, error)

func New(bookingRepo bookingRepository, outboxRepo outboxRepository, bookingMarshaller bookingMarshaller) *service {
	return &service{repo: bookingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

func (b *service) BookTraining(ctx context.Context, booking model.TrainingBooking) error {
	return b.repo.InTx(ctx, func(ctx context.Context) error {
		if err := b.repo.CreateTrainingBooking(ctx, booking); err != nil {
			return err
		}
		event := model.CreateBookingEvent{
			BookingID: booking.TrainingID,
		}

		msg, err := b.bookingMarshaller(event)
		if err != nil {
			return err
		}

		err = b.outboxRepo.CreateOutboxItem(ctx, outbox.OutboxItem{
			Msg:   string(msg),
			Key:   booking.TrainingID,
			Topic: "someTopic",
		})
		if err != nil {
			return err
		}
		return nil

	})
}

func (t *service) ListBookings(ctx context.Context, page page.Page, bookedBy string) ([]model.TrainingBooking, bool, error) {
	models, hasNext, err := t.repo.ListBookings(ctx, page, bookedBy)
	if err != nil {
		return nil, false, err
	}
	return models, hasNext, nil
}

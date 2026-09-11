package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
)

func (b *service) BookTraining(ctx context.Context, booking model.TrainingBooking, withLock bool) error {
	return b.repo.InTx(ctx, func(ctx context.Context) error {
		if err := b.repo.CheckIntersections(ctx, booking); err != nil {
			return err
		}
		if err := b.repo.CreateTrainingBooking(ctx, booking, withLock); err != nil {
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

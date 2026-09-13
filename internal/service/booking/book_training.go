package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

func (b *service) BookTraining(ctx context.Context, booking model.TrainingBooking) error {
	book := model.ListBookingsRequest{
		Page:   page.Page{},
		Filter: model.Filter{BookedFrom: booking.BookFrom, BookedBy: booking.BookedBy, BookedTo: booking.BookTo, WithLock: true},
	}

	return b.repo.InTx(ctx, func(ctx context.Context) error {
		res, err := b.repo.ListBookings(ctx, book)
		if err != nil {
			return err
		}
		if len(res.ModelList) != 0 {
			return model.ErrBookingAlreadyExists
		}
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

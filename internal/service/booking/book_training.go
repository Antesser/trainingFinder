package booking

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

func (b *service) BookTraining(ctx context.Context, booking model.TrainingBooking) error {
	err := b.repo.InTx(ctx, func(ctx context.Context) error {
		res, err := b.repo.ListBookings(ctx, model.ListBookingsRequest{
			Page: page.Page{},
			Filter: model.Filter{
				BookedFrom: &booking.BookFrom,
				BookedTo:   &booking.BookTo,
				WithLock:   true,
			},
		})
		if err != nil {
			return fmt.Errorf("list bookings: %w", err)
		}
		if len(res.ModelList) != 0 {
			return model.ErrBookingAlreadyExists
		}

		booking.Status = model.StatusAwaitingForApproval
		id, err := b.repo.CreateTrainingBooking(ctx, booking)
		if err != nil {
			return fmt.Errorf("create training booking: %w", err)
		}

		if err = b.sendNotification(ctx, model.BookingStatus{BookingID: id, Status: booking.Status}); err != nil {
			return fmt.Errorf("send booking notification: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("book training (in tx): %w", err)
	}
	return nil
}

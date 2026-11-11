package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

func (b *service) BookTraining(ctx context.Context, booking model.TrainingBooking) error {
	return b.repo.InTx(ctx, func(ctx context.Context) error {
		res, err := b.repo.ListBookings(ctx, model.ListBookingsRequest{
			Page: page.Page{},
			Filter: model.Filter{
				BookedFrom: &booking.BookFrom,
				BookedTo:   &booking.BookTo,
				WithLock:   true,
			},
		})
		if err != nil {
			return err
		}
		if len(res.ModelList) != 0 {
			return model.ErrBookingAlreadyExists
		}
		booking.Status = model.StatusAwaitingForApproval
		id, err := b.repo.CreateTrainingBooking(ctx, booking)
		if err != nil {
			return err
		}
		err = b.sendNotification(ctx, model.BookingStatus{BookingID: id, Status: booking.Status})
		if err != nil {
			return err
		}
		return nil

	})
}

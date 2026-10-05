package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
)

func (b *service) ChangeStatusBooking(ctx context.Context, data model.StatusBooking) error {

	err := b.repo.InTx(ctx, func(ctx context.Context) error {
		currentStatus, err := b.repo.GetStatus(ctx, data.BookingID)
		if err != nil {
			return err
		}

		if !model.StatusTransitionAllowed(model.Status(currentStatus), data.Status) {
			return model.ErrBookingStatusTransition
		}

		err = b.repo.ChangeStatus(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

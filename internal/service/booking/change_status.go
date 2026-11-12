package booking

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
)

func (b *service) ChangeBookingStatus(ctx context.Context, data model.BookingStatus) error {
	err := b.repo.InTx(ctx, func(ctx context.Context) error {
		currentStatus, err := b.repo.GetStatus(ctx, data.BookingID)
		if err != nil {
			return fmt.Errorf("get booking status: %w", err)
		}

		if !model.StatusTransitionAllowed(currentStatus, data.Status) {
			return model.ErrBookingStatusTransition
		}

		if err = b.repo.ChangeStatus(ctx, data); err != nil {
			return fmt.Errorf("change booking status: %w", err)
		}

		if err = b.sendNotification(ctx, data); err != nil {
			return fmt.Errorf("send booking status notification: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("change booking status (in tx): %w", err)
	}
	return nil
}

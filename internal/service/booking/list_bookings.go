package booking

import (
	"context"
	"fmt"

	"github.com/Antesser/trainingFinder/internal/model/booking"
)

func (b *service) ListBookings(ctx context.Context, data booking.ListBookingsRequest) (*booking.ListBookingsResponse, error) {
	resp, err := b.repo.ListBookings(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}
	return resp, nil
}

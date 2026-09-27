package booking

import (
	"context"
	"fmt"

	"github.com/Antesser/trainingFinder/internal/model/booking"
)

func (b *service) ListBookings(ctx context.Context, data booking.ListBookingsRequest) (booking.ListBookingsResponse, error) {
	fmt.Println("ddd", data.Filter.BookedBy)
	resp, err := b.repo.ListBookings(ctx, data)
	if err != nil {
		return booking.ListBookingsResponse{}, err
	}
	return resp, nil
}

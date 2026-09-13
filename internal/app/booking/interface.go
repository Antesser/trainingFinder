package booking

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/booking"
)

type bookingService interface {
	BookTraining(ctx context.Context, mod booking.TrainingBooking) error
	ListBookings(ctx context.Context, data booking.ListBookingsRequest) (booking.ListBookingsResponse, error)
}

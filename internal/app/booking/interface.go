package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

type bookingService interface {
	BookTraining(ctx context.Context, mod model.TrainingBooking) error
	ListBookings(ctx context.Context, page page.Page, bookedBy string) ([]model.TrainingBooking, bool, error)
}

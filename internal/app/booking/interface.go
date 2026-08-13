package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/page"
	v1 "github.com/Antesser/trainingFinder/pkg/api/booking/v1"
)

type bookingService interface {
	BookTraining(ctx context.Context, mod model.TrainingBooking) error
	ListBookings(ctx context.Context, page page.Page, filter v1.Filter) ([]model.TrainingBooking, bool, error)
}

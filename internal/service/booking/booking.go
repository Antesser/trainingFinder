package booking

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/page"
)

type service struct {
	repo bookingRepository
}

func New(bookingRepo bookingRepository) *service {
	return &service{repo: bookingRepo}
}
func (t *service) BookTraining(ctx context.Context, booking model.TrainingBooking) error {
	err := t.repo.CreateTrainingBooking(ctx, booking)
	if err != nil {
		return err
	}
	return nil
}

func (t *service) ListBookings(ctx context.Context, page page.Page, bookedBy string) ([]model.TrainingBooking, bool, error) {
	models, hasNext, err := t.repo.ListBookings(ctx, page, bookedBy)
	if err != nil {
		return nil, false, err
	}
	return models, hasNext, nil
}

package booking

import (
	"context"
)

type service struct {
	repo bookingRepository
}

func New(bookingRepo bookingRepository) *service {
	return &service{repo: bookingRepo}
}
func (t *service) BookTraining(ctx context.Context, trainingID, userID string) error {
	err := t.repo.BookTraining(ctx, trainingID, userID)
	if err != nil {
		return err
	}
	return nil
}

package booking

import (
	"context"
)

type bookingService interface {
	BookTraining(ctx context.Context, trainingID, userID string) error
}

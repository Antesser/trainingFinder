package booking

import (
	"context"
)

type bookingRepository interface {
	BookTraining(ctx context.Context, trainingID, userID string) error
}

package training

import (
	"context"
	model "trainingFinder/internal/model/training"
)

type trainingService interface {
	CreateTraining(ctx context.Context, training *model.Training) error
}

package training

import (
	"context"
	model "trainingFinder/internal/model/training"
)

type trainingRepository interface {
	CreateTraining(ctx context.Context, training model.Training) error
	GetTraining(ctx context.Context, id string) (*model.Training, error)
	UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error
	DeleteTraining(ctx context.Context, id string) error
}

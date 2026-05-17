package training

import (
	"context"
	"errors"
	model "trainingFinder/internal/model/training"
)

type trainingService struct {
	repo trainingRepository
}

func New(t trainingRepository) *trainingService {
	return &trainingService{repo: t}
}

func (t *trainingService) CreateTraining(ctx context.Context, trainingModel *model.Training) error {
	if trainingModel == nil {
		return errors.New("training model cannot be nil")
	}
	if err := t.repo.CreateTraining(ctx, *trainingModel); err != nil {
		return err
	}

	return nil
}

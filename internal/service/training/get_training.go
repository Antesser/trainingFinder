package training

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) GetTraining(ctx context.Context, id string) (*model.Training, error) {
	model, err := t.repo.GetTraining(ctx, id)
	if err != nil {
		return model, err
	}
	return model, nil
}

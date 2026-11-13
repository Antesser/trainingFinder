package training

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) GetTraining(ctx context.Context, id string) (*model.Training, error) {
	training, err := t.repo.GetTraining(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get training: %w", err)
	}
	return training, nil
}

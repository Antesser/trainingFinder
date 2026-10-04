package training

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error {
	err := t.repo.UpdateTraining(ctx, updateTraining)
	if err != nil {
		return err
	}
	return nil
}

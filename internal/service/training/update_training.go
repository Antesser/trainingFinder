package training

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error {
	if err := t.repo.UpdateTraining(ctx, updateTraining); err != nil {
		return fmt.Errorf("update training: %w", err)
	}
	return nil
}

package training

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) ListTrainings(ctx context.Context, data model.ListTrainingRequest) (model.ListTrainingResponse, error) {
	resp, err := t.repo.ListTrainings(ctx, data)
	if err != nil {
		return model.ListTrainingResponse{}, fmt.Errorf("list trainings: %w", err)
	}
	return resp, nil
}

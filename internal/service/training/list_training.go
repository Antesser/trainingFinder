package training

import (
	"context"

	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) ListTraining(ctx context.Context, data model.ListTrainingRequest) (model.ListTrainingResponse, error) {
	resp, err := t.repo.ListTrainings(ctx, data)
	if err != nil {
		return model.ListTrainingResponse{}, err
	}
	return resp, nil
}

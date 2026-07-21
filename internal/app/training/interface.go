package training

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/page"
	model "github.com/Antesser/trainingFinder/internal/model/training"
)

type trainingService interface {
	CreateTraining(ctx context.Context, training *model.Training) error
	GetTraining(ctx context.Context, id string) (*model.Training, error)
	UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error
	DeleteTraining(ctx context.Context, id string) error
	ListTraining(ctx context.Context, page page.Page, userID string) ([]model.Training, bool, error)
}

package training

import (
	"context"
	model "trainingFinder/internal/model/training"
)

type trainingRepository interface {
	CreateTraining(ctx context.Context, training *model.Training)
} // вынести декомпозировать Виталий сказал 100% обязательно воскресенье 14,

type trainingService struct {
	repo trainingRepository
}

func New(t trainingRepository) *trainingService {
	return &trainingService{repo: t}
}

func (t *trainingService) CreateTraining(ctx context.Context, trainingModel *model.Training) {
	t.repo.CreateTraining(ctx, trainingModel)
}

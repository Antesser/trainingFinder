package training

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/outbox"
	"github.com/Antesser/trainingFinder/internal/model/page"
	model "github.com/Antesser/trainingFinder/internal/model/training"
	v1 "github.com/Antesser/trainingFinder/pkg/api/training/v1"
)

type trainingMarshaller func(event model.CreateTrainingEvent) ([]byte, error)

type trainingService struct {
	repo               trainingRepository
	outboxRepo         outboxRepository
	trainingMarshaller trainingMarshaller
}

func New(t trainingRepository, trainingMarshaller trainingMarshaller, outboxRepository outboxRepository) *trainingService {
	return &trainingService{repo: t, trainingMarshaller: trainingMarshaller, outboxRepo: outboxRepository}
}

func (t *trainingService) CreateTraining(ctx context.Context, trainingModel *model.Training) error {
	err := t.repo.InTx(ctx, func(ctx context.Context) error {
		if err := t.repo.CreateTraining(ctx, *trainingModel); err != nil {
			return err
		}

		event := model.CreateTrainingEvent{
			TrainingID: trainingModel.ID,
		}

		msg, err := t.trainingMarshaller(event)
		if err != nil {
			return err
		}

		err = t.outboxRepo.CreateOutboxItem(ctx, outbox.OutboxItem{
			Msg:   string(msg),
			Key:   trainingModel.ID,
			Topic: "someTopic",
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (t *trainingService) GetTraining(ctx context.Context, id string) (*model.Training, error) {
	model, err := t.repo.GetTraining(ctx, id)
	if err != nil {
		return model, err
	}
	return model, nil
}
func (t *trainingService) UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error {
	err := t.repo.UpdateTraining(ctx, updateTraining)
	if err != nil {
		return err
	}
	return nil
}
func (t *trainingService) DeleteTraining(ctx context.Context, id string) error {
	err := t.repo.DeleteTraining(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *trainingService) ListTraining(ctx context.Context, page page.Page, filter v1.Filter) ([]model.Training, bool, error) {
	models, hasNext, err := t.repo.ListTrainings(ctx, page, filter)
	if err != nil {
		return nil, false, err
	}
	return models, hasNext, nil
}

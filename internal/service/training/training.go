package training

import (
	"context"
	"errors"
	"trainingFinder/internal/model/outbox"
	model "trainingFinder/internal/model/training"
)

type trainingMarshaller func(event model.CreateTrainingEvent) ([]byte, error)

type trainingService struct {
	repo               trainingRepository
	outboxRepo         outboxRepository
	trainingMarshaller trainingMarshaller
}

func New(t trainingRepository, trainingMarshaller trainingMarshaller) *trainingService {
	return &trainingService{repo: t, trainingMarshaller: trainingMarshaller}
}

func (t *trainingService) CreateTraining(ctx context.Context, trainingModel *model.Training) error {
	if trainingModel == nil {
		return errors.New("training model cannot be nil")
	}

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

	t.outboxRepo.CreateOutboxItem(ctx, outbox.OutboxItem{
		Msg: string(msg),
	})

	return nil
}
func (t *trainingService) GetTraining(ctx context.Context, id string) (*model.Training, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
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
	if id == "" {
		return errors.New("id cannot be empty")
	}
	err := t.repo.DeleteTraining(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

package training

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/outbox"
	model "github.com/Antesser/trainingFinder/internal/model/training"
)

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

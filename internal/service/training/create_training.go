package training

import (
	"context"
	"fmt"

	"github.com/Antesser/trainingFinder/internal/model/outbox"
	model "github.com/Antesser/trainingFinder/internal/model/training"
)

func (t *trainingService) CreateTraining(ctx context.Context, trainingModel *model.Training) error {
	err := t.repo.InTx(ctx, func(ctx context.Context) error {
		if err := t.repo.CreateTraining(ctx, *trainingModel); err != nil {
			return fmt.Errorf("create training: %w", err)
		}

		if err := t.outboxRepo.CreateOutboxItem(ctx, outbox.OutboxItem{
			Msg:   "Later",
			Key:   trainingModel.ID,
			Topic: "someTopic",
		}); err != nil {
			return fmt.Errorf("create outbox item: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("create training (in tx): %w", err)
	}
	return nil
}

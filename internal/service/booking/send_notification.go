package booking

import (
	"context"
	"fmt"

	"github.com/Antesser/trainingFinder/internal/kafka"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
)

// читаем таблицу тренировок и пишем в оутбокс (два метода),
// приватный, в каждой сервисной ручке пихать и ХЕРАЧИТЬ, предварительно обогатив, но ХЕРАЧИТЬ ЖОСКА!!!11
func (b *service) sendNotification(ctx context.Context, trainingID string) error {
	return b.repo.InTx(ctx, func(ctx context.Context) error {
		training, err := b.trainingRepo.GetTraining(ctx, trainingID)
		if err != nil {
			return fmt.Errorf("get training: %w", err)
		}

		payload, err := b.bookingMarshaller(training)
		if err != nil {
			return fmt.Errorf("marshal training: %w", err)
		}

		if err := b.outboxRepo.CreateOutboxItem(ctx, outbox.OutboxItem{
			Msg:   string(payload),
			Topic: kafka.TrainingNotificationsTopic,
			Key:   training.ID,
		}); err != nil {
			return fmt.Errorf("create outbox item: %w", err)
		}

		return nil
	})
}

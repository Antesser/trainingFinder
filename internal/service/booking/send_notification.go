package booking

import (
	"context"
	"fmt"

	"github.com/Antesser/trainingFinder/internal/kafka"
	modelBooking "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/outbox"
)

// читаем таблицу тренировок и пишем в оутбокс (два метода),
// приватный, в каждой сервисной ручке пихать и ХЕРАЧИТЬ, предварительно обогатив, но ХЕРАЧИТЬ ЖОСКА!!!11
func (b *service) sendNotification(ctx context.Context, bookingStatus modelBooking.BookingStatus) error {
	return b.repo.InTx(ctx, func(ctx context.Context) error {

		training, err := b.trainingBookingRepo.GetTrainingByBookingID(ctx, bookingStatus.BookingID)
		if err != nil {
			return fmt.Errorf("get training: %w", err)
		}

		marshallerData := modelBooking.TrainingBookingStatus{TrainingName: training.Name, BookingID: bookingStatus.BookingID, Status: bookingStatus.Status, UserId: training.UserID, Duration: training.Duration}

		payload, err := b.bookingMarshaller(&marshallerData)
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

package outbox

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
)

const (
	limit = 1000
)

type process struct {
	repo     outboxRepository
	producer producer
}

func New(repo outboxRepository,
	producer producer,
) *process {
	return &process{
		repo:     repo,
		producer: producer,
	}
}

func (p *process) ProduceOutboxMessages(ctx context.Context, topic string) error {
	items, err := p.repo.ListOutboxItems(ctx, topic, limit)
	if err != nil {
		return fmt.Errorf("list outbox messages: %w", err)
	}

	if len(items) == 0 {
		return nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(*items[j].CreatedAt)
	})

	deleteIDs := make([]uuid.UUID, 0, len(items))

	for _, item := range items {
		if err := p.producer.Produce(item.Topic, item.Msg, item.Key); err != nil {
			// log.Errorf("produce outbox message %s: %v", item.MessageValue, err)
			// Нужно выйти из цикла, чтобы сохранить порядок отправки
			break
		}

		deleteIDs = append(deleteIDs, item.ID)
	}

	// Если не получится удалить сообщение из outbox-таблицы,
	// попробуем в следующей итерации.
	// Гарантия доставки – at least once.
	if err := p.repo.DeleteOutboxItem(ctx, deleteIDs); err != nil {
		return fmt.Errorf("delete outbox items: %w", err)
	}
	return nil
}

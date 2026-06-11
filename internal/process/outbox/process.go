package outbox

import (
	"context"
	"fmt"
	"sort"
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
		return items[i].ID < items[j].ID
	})

	deleteIDs := make([]int64, 0, len(items))

	for _, item := range items {
		if err := p.producer.Produce(ctx, item.Topic, []byte(item.Msg), item.Key); err != nil {
			// log.Errorf("produce outbox message %s: %v", item.MessageValue, err)
			// Нужно выйти из цикла, чтобы сохранить порядок отправки
			break
		}

		deleteIDs = append(deleteIDs, item.ID)
	}

	// Если не получится удалить сообщение из outbox-таблицы,
	// попробуем в следующей итерации.
	// Гарантия доставки – at least once.
	return p.repo.DeleteOutboxItem(ctx, deleteIDs)
}

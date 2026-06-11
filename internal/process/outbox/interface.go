package outbox

import (
	"context"
	"trainingFinder/internal/model/outbox"
)

type outboxRepository interface {
	DeleteOutboxItem(ctx context.Context, id []int64) error
	ListOutboxItems(ctx context.Context, topic string, limit uint64) ([]outbox.OutboxItem, error)
}

type producer interface {
	Produce(ctx context.Context, topic string, msg []byte, key string) error
}

package outbox

import (
	"context"

	"github.com/Antesser/trainingFinder/internal/model/outbox"
	"github.com/google/uuid"
)

type outboxRepository interface {
	DeleteOutboxItem(ctx context.Context, id []uuid.UUID) error
	ListOutboxItems(ctx context.Context, topic string, limit uint64) ([]outbox.OutboxItem, error)
}

type producer interface {
	Produce(topic, msg, aggregateKey string) error
}

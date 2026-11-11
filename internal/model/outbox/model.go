package outbox

import (
	"time"

	"github.com/google/uuid"
)

type OutboxItem struct {
	ID        uuid.UUID
	Msg       string
	Topic     string
	Key       string // записать сюда trASININGididididid
	CreatedAt *time.Time
}

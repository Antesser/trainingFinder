package training

import (
	"context"
	"trainingFinder/internal/model/outbox"
	model "trainingFinder/internal/model/training"

	"github.com/jackc/pgx/v5"
)

type trainingRepository interface {
	CreateTraining(ctx context.Context, training model.Training) error
	GetTraining(ctx context.Context, id string) (*model.Training, error)
	UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error
	DeleteTraining(ctx context.Context, id string) error
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	InTxWithIsoLevel(ctx context.Context, isoLevel pgx.TxIsoLevel, fn func(ctx context.Context) error) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
	DeleteOutboxItem(ctx context.Context, id []string) error
}

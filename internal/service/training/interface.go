package training

import (
	"context"
	"trainingFinder/internal/model/outbox"
	model "trainingFinder/internal/model/training"

	"github.com/jackc/pgx/v5"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name trainingRepository --exported
type trainingRepository interface {
	CreateTraining(ctx context.Context, training model.Training) error
	GetTraining(ctx context.Context, id string) (*model.Training, error)
	UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error
	DeleteTraining(ctx context.Context, id string) error
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	InTxWithIsoLevel(ctx context.Context, isoLevel pgx.TxIsoLevel, fn func(ctx context.Context) error) error
	ListTrainings(ctx context.Context) ([]model.Training, error)
	BookTraining(ctx context.Context, trainingID, userID string) error
}

type outboxRepository interface {
	CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error
}

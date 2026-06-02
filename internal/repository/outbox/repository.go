package outbox

import (
	"context"
	"fmt"
	"trainingFinder/internal/model/outbox"
	model "trainingFinder/internal/model/training"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/pgxtransactor"
)

type repository struct {
	pool *pgxtransactor.Pool
}

type outboxItem struct {
	ID           string `db:"id"`
	Topic        string `db:"topic"`
	Key          string `db:"key"`
	MessageValue string `db:"message_value"`
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) CreateOutboxItem(ctx context.Context, item outbox.OutboxItem) error {
	qb := sq.Insert("outbox").
		Columns("message_value", "channel", "key").
		Values(item.Msg, item.Topic, item.Key).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {

		return err
	}

	if _, err = r.pool.Querier(ctx).Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (r *repository) ListOutboxItems(ctx context.Context, limit uint64) ([]outbox.OutboxItem, error) {
	qb := sq.Select("message_value", "channel", "key").
		From("outbox").
		Limit(limit).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var items []outboxItem
	err = pgxscan.Select(ctx, r.pool.Querier(ctx), &items, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	result := make([]outbox.OutboxItem, 0, len(items))
	for _, i := range items {
		result = append(result, outbox.OutboxItem{
			Msg:   i.MessageValue,
			Topic: i.Topic,
			Key:   i.Key,
		})
	}

	return result, nil
}

func (r *repository) DeleteOutboxItem(ctx context.Context, id []string) error {
	qb := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	tags, err := r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tags.RowsAffected() == 0 {
		return model.ErrTrainingNotFound
	}
	return nil
}

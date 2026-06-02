package outbox

import (
	"context"
	"fmt"
	"trainingFinder/internal/model/outbox"
	model "trainingFinder/internal/model/training"

	sq "github.com/Masterminds/squirrel"
	"github.com/golangmonster/pgxtransactor"
)

type repository struct {
	pool *pgxtransactor.Pool
}

type outboxItem struct {
	ID             string    `db:"id"`
	key         string    `db:"key"`
	message_value      string `db:"message_value"`
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
		Limit("limit").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

var items []outboxItem
err = pgxscan.Select(ctx, r.pool, &items, query, args...){
		var item outbox.OutboxItem
		if err := rows.Scan(&item.Msg, &item.Topic, &item.Key); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
	

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *repository) DeleteOutboxItem(ctx context.Context, id string) error {
	qb := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if tags, err := r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}
	if tags.RowsAffected() == 0 {
			return model.ErrTrainingNotFound
		}
	return nil
}

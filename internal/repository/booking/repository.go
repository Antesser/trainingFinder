package booking

import (
	"context"
	model "trainingFinder/internal/model/booking"

	sq "github.com/Masterminds/squirrel"
	"github.com/golangmonster/pgxtransactor"
)

type repository struct {
	pool *pgxtransactor.Pool
	pgxtransactor.Transactor
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool, Transactor: pool}
}
func (r *repository) BookTraining(ctx context.Context, trainingID, userID string) error {
	qb := sq.Update("training").
		Where(sq.Eq{"id": trainingID})

	qb = qb.Set("booked_by", userID)

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tags, err := r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tags.RowsAffected() == 0 {
		return model.ErrBookingNotFound
	}
	return nil
}

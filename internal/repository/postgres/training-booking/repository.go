package training_booking

import (
	"context"
	"database/sql"
	"errors"

	model "github.com/Antesser/trainingFinder/internal/model/training"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/pgxtransactor"
)

type repository struct {
	pgxtransactor.Transactor

	pool *pgxtransactor.Pool
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool, Transactor: pool}
}

func (r *repository) GetTrainingByBookingID(ctx context.Context, bookingID string) (*model.Training, error) {
	qb := sq.Select(
		"id",
		"trainer_id", "user_id", "started_at", "ended_at", "additional_info", "training_name",
	).From("training as t").
		LeftJoin("training_booking as tb on tb.training_id=t.id").
		Where(sq.Eq{"tb.id": bookingID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var t training
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &t, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTrainingNotFound
		}
	}

	return &model.Training{ID: t.ID, UserID: t.UserID, TrainerID: t.TrainerID, StartedAt: t.StartedAt, EndedAt: t.EndedAt}, nil
}

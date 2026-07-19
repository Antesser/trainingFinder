package training

import (
	"context"
	"database/sql"
	"errors"
	"time"
	model "trainingFinder/internal/model/training"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/pgxtransactor"
)

type training struct {
	ID             string    `db:"id"`
	TrainerID      string    `db:"trainer_id"`
	UserID         string    `db:"user_id"`
	StartedAt      time.Time `db:"started_at"`
	EndedAt        time.Time `db:"ended_at"`
	AdditionalInfo string    `db:"additional_info"`
}
type repository struct {
	pool *pgxtransactor.Pool
	pgxtransactor.Transactor
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool, Transactor: pool}
}

func (r *repository) CreateTraining(ctx context.Context, training model.Training) error {
	qb := sq.Insert("training").
		Columns("id", "trainer_id", "user_id", "started_at", "ended_at", "additional_info").
		Values(training.ID, training.TrainerID, training.UserID, training.StartedAt, training.EndedAt, training.AdditionalInfo).
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
func (r *repository) GetTraining(ctx context.Context, id string) (*model.Training, error) {
	qb := sq.Select(
		"id",
		"trainer_id", "user_id", "started_at", "ended_at", "additional_info",
	).From("training").
		Where(sq.Eq{"id": id}).
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
func (r *repository) UpdateTraining(ctx context.Context, updateTraining model.UpdateTrainingRequest) error {
	qb := sq.Update("training").
		Where(sq.Eq{"id": updateTraining.UserID})
	if updateTraining.AdditionalInfo != nil {
		qb = qb.Set("additional_info", updateTraining.AdditionalInfo)
	}
	if updateTraining.TrainerID != nil {
		qb = qb.Set("trainer_id", updateTraining.TrainerID)
	}
	if updateTraining.UserID != nil {
		qb = qb.Set("user_id", updateTraining.UserID)
	}
	if updateTraining.StartedAt != nil {
		qb = qb.Set("started_at", updateTraining.StartedAt)
	}
	if updateTraining.EndedAt != nil {
		qb = qb.Set("ended_at", updateTraining.EndedAt)
	}

	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
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
func (r *repository) DeleteTraining(ctx context.Context, id string) error {
	qb := sq.Delete("training").
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

func (r *repository) ListTrainings(ctx context.Context, pageLimit, offset uint64, userID string) ([]model.Training, error) {
	qb := sq.Select(
		"id",
		"trainer_id", "user_id", "started_at", "ended_at", "additional_info",
	).From("training").
		Where(sq.Eq{"user_id": userID}).
		Limit(pageLimit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var rows []training
	if err := pgxscan.Select(ctx, r.pool.Querier(ctx), &rows, query, args...); err != nil {
		return nil, err
	}

	out := make([]model.Training, 0, len(rows))
	for i, t := range rows {
		out[i] = model.Training{ID: t.ID, TrainerID: t.TrainerID, UserID: t.UserID, StartedAt: t.StartedAt, EndedAt: t.EndedAt}
	}

	return out, nil
}

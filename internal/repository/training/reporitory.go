package training

import (
	"context"
	model "trainingFinder/internal/model/training"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func (r *repository) CreateTraining(ctx context.Context, training model.Training) error {
	qb := sq.Insert("training").
		Columns("trainer_id", "user_id", "started_at", "ended_at", "additional_info").
		Values(training.TrainerID, training.UserID, training.StartedAt, training.EndedAt, training.AdditionalInfo).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

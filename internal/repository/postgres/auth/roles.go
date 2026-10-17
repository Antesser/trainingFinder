package auth

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/training"
	"github.com/Antesser/trainingFinder/internal/utils"
	sq "github.com/Masterminds/squirrel"
)

func (r *repository) CreateRole(ctx context.Context, role string) error {
	qb := sq.Insert("roles").
		Columns("role").
		Values(role).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		if utils.IsUniqueViolation(err) {
			return model.ErrAlreadyExists
		}
		return fmt.Errorf("execute insert: %w", err)
	}
	return nil
}

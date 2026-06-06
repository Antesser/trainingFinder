package auth

import (
	"context"
	model "trainingFinder/internal/model/auth"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) SaveRefreshToken(ctx context.Context, session model.Sessions) error {
	qb := sq.Insert("sessions").
		Columns("refresh_token", "user_id", "is_active", "created_at", "expires_at").
		Values(session.Token, session.UserID, session.Active, session.CreatedAt, session.ExpiresAt).
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

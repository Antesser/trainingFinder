package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	model "trainingFinder/internal/model/auth"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type sessions struct {
	Token     uuid.UUID `db:"refresh_token"`
	UserID    string    `db:"user_id"`
	Active    bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (r *repository) GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*model.Sessions, error) { // в транзакцию вставка в таблицу сессий
	qb := sq.Select("refresh_token", "user_id", "is_active", "created_at", "expires_at").
		From("session").
		Where(sq.And{sq.Eq{"refresh_token": refreshToken}, sq.Eq{"is_active": true}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	var i sessions
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &i, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute insert: %w", err)
	}

	return &model.Sessions{Token: i.Token, UserID: i.UserID, Active: i.Active,
		CreatedAt: i.CreatedAt, ExpiresAt: i.ExpiresAt}, nil
}

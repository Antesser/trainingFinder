package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/auth"

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
	RoleID    string    `db:"role_id"`
}

func (r *repository) GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*model.Session, error) { // в транзакцию вставка в таблицу сессий
	qb := sq.Select("s.refresh_token", "s.user_id", "s.is_active", "s.created_at", "s.expires_at", "u.role_id").
		From("session s").
		Join("users u ON u.id = s.user_id").
		Where(sq.And{sq.Eq{"s.refresh_token": refreshToken}, sq.Eq{"s.is_active": true}}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	var i sessions
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &i, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrSessionNotFound
		}
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return &model.Session{Token: i.Token, UserID: i.UserID, Active: i.Active,
		CreatedAt: i.CreatedAt, ExpiresAt: i.ExpiresAt, RoleID: i.RoleID}, nil
}

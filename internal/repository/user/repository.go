package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	model "trainingFinder/internal/model/user"

	sq "github.com/Masterminds/squirrel"
	"github.com/golangmonster/pgxtransactor"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	pool *pgxtransactor.Pool
}
type training struct {
	ID             string    `db:"id"`
	username      string    `db:"username"`
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) GetUserByID(ctx context.Context, id string) (string, error) {
	qb := sq.Select(
		"id",
		"username",
	).From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return "", err
	}

	var u user
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &u, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", model.ErrUserNotFound
		}
		return "", fmt.Errorf("execute insert: %w", err)
	}

	return username, nil
}

func (r *repository) UpdateUser(ctx context.Context, userID, username string) error {
	qb := sq.Update("users").
		Where(sq.Eq{"id": userID}).
		Set("username", username).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return err
	}

	returnData, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if returnData.RowsAffected() == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	qb := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	if tags, err := r.pool.Exec(ctx, query, args...); err != nil {

		if tags.RowsAffected() == 0 {
			return model.ErrUserNotFound
		}
		return err
	}
	return nil
}

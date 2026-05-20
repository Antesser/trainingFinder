package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) GetUser(ctx context.Context, id string) (string, error) {
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

	var username string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&username)
	if err != nil {
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

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return err
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

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

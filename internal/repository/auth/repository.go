package auth

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/golangmonster/pgxtransactor"
)

const AlreadyExists = "23505"

type repository struct {
	pool *pgxtransactor.Pool
}
type user struct {
	ID string `db:"id"`
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SignUp(ctx context.Context, hash []byte, id, login string) (string, error) {
	qb := sq.Insert("users").
		Columns("id", "username", "password").
		Values(id, login, hash).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return "", err
	}

	var userID string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("execute insert: %w", err)
	}

	return userID, nil
}

package auth

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/pgxtransactor"
)

const AlreadyExists = "23505"

type repository struct {
	pool *pgxtransactor.Pool
}
type user struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
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

	var u user
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &u, query, args...)
	if err != nil {
		return "", fmt.Errorf("execute insert: %w", err)
	}

	return u.ID, nil
}

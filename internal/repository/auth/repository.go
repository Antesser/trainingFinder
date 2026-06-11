package auth

import (
	"context"
	"fmt"
	"log"
	model "trainingFinder/internal/model/user"

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
	Login    string `db:"login"`
	Password string `db:"password"`
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SignUp(ctx context.Context, hash []byte, id, login string) (string, error) {
	err := r.CheckUserExistence(ctx, login)
	if err != nil {
		return "", err
	}
	qb := sq.Insert("users").
		Columns("id", "login", "password").
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
func (r *repository) CheckUserExistence(ctx context.Context, login string) error {
	checkQuery := sq.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where("login = ?", login).
		Suffix(")").
		PlaceholderFormat(sq.Dollar)

	checkSQL, checkArgs, err := checkQuery.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return err
	}

	var exists bool
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &exists, checkSQL, checkArgs...)
	if err != nil {
		return fmt.Errorf("check login existence: %w", err)
	}

	if exists {
		return model.ErrLoginAlreadyExists
	}
	return nil
}

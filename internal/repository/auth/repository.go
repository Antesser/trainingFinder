package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	model "trainingFinder/internal/model/training"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const AlreadyExists = "23505"

type repository struct {
	pool *pgxpool.Pool
}
type user struct {
	ID string `db:"id"`
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) SignUp(ctx context.Context, hash, login, password string) (string, error) {
	qb := sq.Insert("users").
		Columns("id", "username", "password").
		Values(hash, login, password).
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
func (r *repository) GetUserByIDAuthInfo(ctx context.Context, hash, login, password string) (string, error) { // в транзакцию вставка в таблицу сессий
	qb := sq.Insert("users").
		Columns("id", "username", "password").
		Values(hash, login, password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		// извлекаем конкретный тип из ошибки, проходим по цепочке ошибок и сравниваем, если то, что нужно - AlreadyExists, то возвращаем нашу созданную ошибку
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == AlreadyExists {
			return "", model.ErrAlreadyExists
		}
		return "", err
	}

	var userID string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("execute insert: %w", err)
	}

	return userID, nil
}

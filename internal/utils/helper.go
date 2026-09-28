package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func TruncateForHasNext[T any](slice []T, limit int) ([]T, bool) {
	hasNext := len(slice) > limit
	if hasNext {
		slice = slice[:limit]
	}
	return slice, hasNext
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

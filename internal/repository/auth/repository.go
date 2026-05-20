package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}
type user struct {
	ID string `db:"id"`
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
func generateRandomHash(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func (r *repository) SignUp(ctx context.Context, login, password string) (string, string, error) {
	hash, _ := generateRandomHash(16)
	qb := sq.Insert("users").
		Columns("id", "username", "password").
		Values(hash, login, password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return "", "", err
	}

	var userID string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return "", "", fmt.Errorf("execute insert: %w", err)
	}

	accessToken, refreshToken := "", ""

	return accessToken, refreshToken, nil
}

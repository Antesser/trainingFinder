package auth

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
type user struct {
	ID string `db:"id"`
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
func (r *repository) SignUp(ctx context.Context, login, password string) (string, string, error) {
	qb := sq.Insert("users").
		Columns("name", "password").
		Values(login, password).
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

// func (r *repository) SignIn(username string, userPassword string) (string, error) {
// 	hash, err := password.Hash(userPassword)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	qb := squirrel.Select("id").From("table").Where(
// 		squirrel.And{
// 			squirrel.Eq{"username": username},
// 			squirrel.Or{
// 				squirrel.Neq{"username": nil},
// 				squirrel.Eq{"username": username},
// 			},
// 		})

// 	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()

// 	pool.Exec(sql, args...) // - Выполнение запроса UPDATE/DELETE если нет RETURNING

// 	var user user
// 	// Чтение - SELECT, есть RETURNING
// 	// Несколько строк
// 	pgxscan.Select(ctx, r.c, &user, sql, args...)

//		// Только одна строка
//		pgxscan.Get(ctx, r.c, &user, sql, args...)
//	}

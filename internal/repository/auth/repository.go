package auth

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vzglad-smerti/password_hash"
)

type repository struct {
	pool *pgxpool.Pool
}

type user struct {
	ID int `db:"id"`
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
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
func (r *repository) CreateUser(username string, userpassword string) error {
	hashPassword, err := password.Hash(userpassword)
	if err != nil {
		log.Print(err)
	}

	qb := sq.Insert("users").
		Columns("name", "password").
		Values(username, hashPassword).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return err
	}

	log.Print("SQL:", query)
	log.Print("Args:", args)

	return nil
}

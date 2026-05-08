package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type repository struct {
	c pgxpool.Conn
}

type user struct {
	ID `db:"id"`
}

func (r *repository) SignIn(username string, password string) (string, error) {
	qb := squirrel.Select("id").From("table").Where(
		squirrel.And{
			squirrel.Eq{"username": username},
			squirrel.Or{
				squirrel.Neq{"username": nil},
				squirrel.Eq{"username": username},
			},
		})

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()

	c.Exec(sql, args...) // - Выполнение запроса UPDATE/DELETE если нет RETURNING

	var user user
	// Чтение - SELECT, есть RETURNING
	// Несколько строк
	pgxscan.Select(ctx, r.c, &user, sql, args...)

	// Только одна строка
	pgxscan.Get(ctx, r.c, &user, sql, args...)
}

package auth

import (
	"context"
	"errors"
	"fmt"
	model "trainingFinder/internal/model/auth"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type userAuthRepo struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r *repository) GetUserByIDAuthInfo(ctx context.Context, login string) (model.UserAuthInfo, error) { // в транзакцию вставка в таблицу сессий
	qb := sq.Select("users").
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		// извлекаем конкретный тип из ошибки, проходим по цепочке ошибок и сравниваем, если то, что нужно - AlreadyExists, то возвращаем нашу созданную ошибку
		if errors.Is(err, pgx.ErrNoRows) {
			return model.UserAuthInfo{}, model.ErrNotFound
		}
		return model.UserAuthInfo{}, err
	}
	returnModel := userAuthRepo{}
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &returnModel, query, args...)
	if err != nil {
		return model.UserAuthInfo{}, fmt.Errorf("execute insert: %w", err)
	}

	return model.UserAuthInfo{ID: returnModel.ID, Login: returnModel.Login, Password: returnModel.Password}, nil
}

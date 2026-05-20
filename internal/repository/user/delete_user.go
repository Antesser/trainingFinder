package user

//
//import (
//	"context"
//
//	sq "github.com/Masterminds/squirrel"
//)
//
//func (r *repository) DeleteUser(ctx context.Context, id string) error {
//	qb := sq.Delete("users").
//		Where(sq.Eq{"id": id}).
//		PlaceholderFormat(sq.Dollar)
//
//	query, args, err := qb.ToSql()
//	if err != nil {
//		return err
//	}
//
//	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
//		return err
//	}
//	return nil
//}

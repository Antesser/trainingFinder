package auth

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *repository) CreateRole(ctx context.Context, role string) error {
	qb := sq.Insert("roles").
		Columns("role").
		Values(role).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		log.Print("Got an error:", err)
		return err
	}

	var userRole roleStruct
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &userRole, query, args...)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}

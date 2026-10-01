package booking

import (
	"context"
	"fmt"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/utils"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/samber/lo"

	sq "github.com/Masterminds/squirrel"
	"github.com/golangmonster/pgxtransactor"
)

type repository struct {
	pool *pgxtransactor.Pool
	pgxtransactor.Transactor
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool, Transactor: pool}
}

func (r *repository) CreateTrainingBooking(ctx context.Context, booking model.TrainingBooking) error {
	qb := sq.Insert("training_booking").
		Columns("training_id", "booked_by", "book_to", "book_from", "created_at").
		Values(booking.TrainingID, booking.UserID, booking.BookTo, booking.BookFrom, time.Now()).
		PlaceholderFormat(sq.Dollar)
	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ListBookings(ctx context.Context, data model.ListBookingsRequest) (model.ListBookingsResponse, error) {
	limit := int(data.Page.Limit)
	fmt.Println("data", data.Filter.BookedBy)
	qb := sq.Select(
		"id",
		"training_id",
		"booked_by",
		"created_at",
		"status",
		"book_from",
		"book_to",
	).From("training_booking").
		Limit(data.Page.Limit + 1).
		Offset(data.Page.Offset).
		PlaceholderFormat(sq.Dollar)

	if data.Filter.WithLock {
		qb = qb.Suffix("FOR UPDATE")
	}
	if data.Filter.BookedFrom != nil {
		qb = qb.Where(sq.Eq{"book_from": *data.Filter.BookedFrom})
	}
	if data.Filter.BookedBy != nil {
		qb = qb.Where(sq.Eq{"booked_by": *data.Filter.BookedBy})
	}
	if data.Filter.TrainerID != nil {
		qb = qb.Where(sq.Eq{"training_id": *data.Filter.TrainerID})
	}
	if data.Filter.BookedTo != nil {
		qb = qb.Where(sq.Eq{"book_to": *data.Filter.BookedTo})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return model.ListBookingsResponse{}, err
	}

	var rows []booking
	if err := pgxscan.Select(ctx, r.pool.Querier(ctx), &rows, query, args...); err != nil {
		return model.ListBookingsResponse{}, err
	}
	rows, hasNext := utils.TruncateForHasNext(rows, limit)

	out := lo.Map(rows, func(b booking, _ int) model.TrainingBooking {
		return model.TrainingBooking{
			ID:         b.ID,
			TrainingID: b.TrainingID,
			BookFrom:   b.BookFrom,
			BookTo:     b.BookTo,
			UserID:     b.UserID,
			Status:     b.Status,
			CreatedAt:  b.CreatedAt,
		}
	})
	fmt.Println("SQL:", query)
	fmt.Println("ARGS:", args)
	fmt.Println("rows:", len(rows))
	return model.ListBookingsResponse{ModelList: out, HasNext: hasNext}, nil
}

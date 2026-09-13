package booking

import (
	"context"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/utils"
	"github.com/georgysavva/scany/v2/pgxscan"

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

type booking struct {
	ID         string    `db:"id"`
	TrainingID string    `db:"training_id"`
	UserID     string    `db:"booked_by"`
	CreatedAt  time.Time `db:"created_at"`
	BookFrom   time.Time `db:"book_from"`
	BookTo     time.Time `db:"book_to"`
}

func (r *repository) CreateTrainingBooking(ctx context.Context, booking model.TrainingBooking) error {
	qb := sq.Insert("training_booking").
		Columns("training_id", "booked_by", "book_to", "book_from").
		Values(booking.TrainingID, booking.UserID, booking.BookTo, booking.BookFrom).
		PlaceholderFormat(sq.Dollar)
	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tags, err := r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tags.RowsAffected() == 0 {
		return model.ErrBookingNotFound
	}
	return nil
}

func (r *repository) ListBookings(ctx context.Context, data model.ListBookingsRequest) (model.ListBookingsResponse, error) {
	limit := int(data.Page.Limit)
	qb := sq.Select(
		"id",
		"training_id",
		"booked_by",
		"created_at",
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
		qb = qb.Where("book_from = ?", data.Filter.BookedFrom)
	}
	if data.Filter.BookedBy != nil {
		qb = qb.Where("booked_by = ?", data.Filter.BookedBy)
	}
	if data.Filter.BookedTo != nil {
		qb = qb.Where("book_to = ?", data.Filter.BookedTo)
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
	out := make([]model.TrainingBooking, 0, len(rows))
	for i, b := range rows {
		out[i] = model.TrainingBooking{TrainingID: b.TrainingID, BookFrom: &b.BookFrom, BookTo: &b.BookTo, UserID: b.UserID}
	}

	return model.ListBookingsResponse{ModelList: out, HasNext: hasNext}, nil
}

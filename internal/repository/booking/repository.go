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

func (r *repository) CheckIntersections(ctx context.Context, booking model.TrainingBooking) error {
	sel := sq.Select(
		"id",
	).From("training_booking").
		Where(sq.Eq{"training_id": booking.TrainingID}).
		Where(sq.Eq{"booked_from": booking.BookFrom}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sel.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tags, err := r.pool.Querier(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tags.RowsAffected() != 0 {
		return model.ErrBookingAlreadyExists
	}
	return nil
}

func (r *repository) CreateTrainingBooking(ctx context.Context, booking model.TrainingBooking, withLock bool) error {
	qb := sq.Insert("training_booking").
		Columns("training_id", "booked_by", "book_to", "book_from").
		Values(booking.TrainingID, booking.UserID, booking.BookTo, booking.BookFrom).
		PlaceholderFormat(sq.Dollar)
	if withLock {
		qb.Suffix("FOR UPDATE")
	}
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
		"training_id", "booked_by", "created_at", "book_from", "book_to",
	).From("training_booking").
		Where(sq.Eq{"booked_by": data.BookedBy}).
		Limit(data.Page.Limit + 1).
		Offset(data.Page.Offset).
		PlaceholderFormat(sq.Dollar)

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
		out[i] = model.TrainingBooking{TrainingID: b.TrainingID, BookFrom: b.BookFrom, BookTo: b.BookTo, UserID: b.UserID}
	}

	return model.ListBookingsResponse{ModelList: out, HasNext: hasNext}, nil
}

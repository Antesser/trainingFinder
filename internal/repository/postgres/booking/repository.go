package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/pgxtransactor"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

type repository struct {
	pgxtransactor.Transactor

	pool *pgxtransactor.Pool
}

func New(pool *pgxtransactor.Pool) *repository {
	return &repository{pool: pool, Transactor: pool}
}

func (r *repository) CreateTrainingBooking(ctx context.Context, booking model.TrainingBooking) (string, error) {
	qb := sq.Insert("training_booking").
		Columns("training_id", "booked_by", "book_to", "book_from", "created_at", "status").
		Values(booking.TrainingID, booking.UserID, booking.BookTo, booking.BookFrom, time.Now(), booking.Status).
		PlaceholderFormat(sq.Dollar).
		Suffix("returning id")
	query, args, err := qb.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", err
	}
	var id string
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &id, query, args...)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *repository) ListBookings(ctx context.Context, data model.ListBookingsRequest) (*model.ListBookingsResponse, error) {
	limit := data.Page.Limit
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
		return nil, err
	}

	var rows []booking
	if err := pgxscan.Select(ctx, r.pool.Querier(ctx), &rows, query, args...); err != nil {
		return nil, err
	}
	const maxInt = int(^uint(0) >> 1)

	if limit > uint64(maxInt) {
		limit = uint64(maxInt)
	}
	rows, hasNext := utils.TruncateForHasNext(rows, int(limit))

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

	return &model.ListBookingsResponse{ModelList: out, HasNext: hasNext}, nil
}

func (r *repository) GetStatus(ctx context.Context, bookingID string) (model.Status, error) {
	qb := sq.Select("status").
		From("training_booking").
		Where(sq.Eq{"id": bookingID}).
		PlaceholderFormat(sq.Dollar).
		Suffix("FOR UPDATE")

	query, args, err := qb.ToSql()
	if err != nil {
		return "", fmt.Errorf("build query: %w", err)
	}

	var status model.Status
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &status, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", model.ErrBookingNotFound
		}

		return "", fmt.Errorf("query status: %w", err)
	}

	return status, nil
}

func (r *repository) ChangeStatus(ctx context.Context, data model.BookingStatus) error {
	qb := sq.Update("training_booking").
		Set("status", data.Status).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": data.BookingID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	var out model.BookingStatus
	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &out, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrBookingNotFound
		}

		return fmt.Errorf("execute update: %w", err)
	}

	return nil
}

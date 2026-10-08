package booking

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrBookingNotFound = errors.New("booking not found")
var ErrBookingAlreadyExists = errors.New("booking already exists")
var ErrBookingStatusTransition = errors.New("not allowed to change to that status")
var ErrBookingNoRows = errors.New("no rows in result set")

type TrainingBooking struct {
	ID         uuid.UUID
	TrainingID string
	UserID     string
	BookFrom   time.Time
	BookTo     time.Time
	Status     string
	BookedBy   string
	CreatedAt  time.Time
}
type CreateBookingEvent struct {
	BookingID string
}

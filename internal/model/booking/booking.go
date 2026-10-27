package booking

import (
	"errors"
	"slices"
	"time"
	"uuid"
)

type (
	Status string
)

const (
	StatusUnspecified         Status = "UNSPECIFIED"
	StatusAwaitingForApproval Status = "AWAITING_FOR_APPROVAL"
	StatusCancelled           Status = "CANCELLED"
	StatusApproved            Status = "APPROVED"
	StatusInProgress          Status = "IN_PROGRESS"
	StatusCompleted           Status = "COMPLETED"
	StatusDeclined            Status = "DECLINED"
)

var (
	allowedStatusTransitions = map[Status][]Status{
		StatusAwaitingForApproval: {StatusApproved, StatusDeclined},
		StatusApproved:            {StatusInProgress, StatusCancelled},
		StatusInProgress:          {StatusCompleted},
		StatusCancelled:           {},
		StatusCompleted:           {},
		StatusDeclined:            {},
	}
)

func StatusTransitionAllowed(from, to Status) bool {
	return slices.Contains(allowedStatusTransitions[from], to)
}

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

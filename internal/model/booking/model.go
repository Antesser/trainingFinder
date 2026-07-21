package booking

import (
	"errors"
	"time"
)

var ErrBookingNotFound = errors.New("booking not found")
var ErrBookingAlreadyExists = errors.New("booking already exists")

type TrainingBooking struct {
	TrainingID string
	UserID     string
	BookFrom   time.Time
	BookTo     time.Time
}

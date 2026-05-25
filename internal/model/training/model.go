package training

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Training struct {
	ID             string
	TrainerID      string
	UserID         string
	StartedAt      time.Time
	EndedAt        time.Time
	AdditionalInfo string
}

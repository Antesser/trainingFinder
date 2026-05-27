package training

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("user already exists")

type Training struct {
	ID             string
	TrainerID      string
	UserID         string
	StartedAt      time.Time
	EndedAt        time.Time
	AdditionalInfo string
}

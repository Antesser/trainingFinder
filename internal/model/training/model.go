package training

import (
	"errors"
	"time"
)

var ErrTrainingNotFound = errors.New("training not found")
var ErrAlreadyExists = errors.New("user already exists")

type Training struct {
	ID             string
	TrainerID      string
	UserID         string
	StartedAt      time.Time
	EndedAt        time.Time
	AdditionalInfo string
}

type CreateTrainingEvent struct {
	TrainingID string
}

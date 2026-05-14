package training

import "time"

type Training struct {
	TrainerID      string
	UserID         string
	StartedAt      time.Time
	EndedAt        time.Time
	AdditionalInfo string
}

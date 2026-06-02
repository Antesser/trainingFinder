package training

import "time"

type UpdateTrainingRequest struct {
	ID             string
	TrainerID      *string
	UserID         *string
	StartedAt      *time.Time
	EndedAt        *time.Time
	AdditionalInfo *string
}

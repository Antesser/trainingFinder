package training

import "time"

type training struct {
	ID             string    `db:"id"`
	TrainerID      string    `db:"trainer_id"`
	UserID         string    `db:"user_id"`
	StartedAt      time.Time `db:"started_at"`
	EndedAt        time.Time `db:"ended_at"`
	AdditionalInfo string    `db:"additional_info"`
	Duration       int64     `db:"duration"`
}

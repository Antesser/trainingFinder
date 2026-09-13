package booking

import "time"

type booking struct {
	ID         string    `db:"id"`
	TrainingID string    `db:"training_id"`
	UserID     string    `db:"booked_by"`
	CreatedAt  time.Time `db:"created_at"`
	BookFrom   time.Time `db:"book_from"`
	BookTo     time.Time `db:"book_to"`
	Status     string    `db:"status"`
}

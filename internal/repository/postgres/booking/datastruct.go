package booking

import (
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/google/uuid"
)

type booking struct {
	ID         uuid.UUID    `db:"id"`
	TrainingID string       `db:"training_id"`
	UserID     string       `db:"booked_by"`
	CreatedAt  time.Time    `db:"created_at"`
	BookFrom   time.Time    `db:"book_from"`
	BookTo     time.Time    `db:"book_to"`
	Status     model.Status `db:"status"`
}

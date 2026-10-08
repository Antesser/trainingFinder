package booking

import (
	"time"

	"github.com/Antesser/trainingFinder/internal/model/page"
)

type ListBookingsRequest struct {
	Page   page.Page
	Filter Filter
}
type ListBookingsResponse struct {
	ModelList []TrainingBooking
	HasNext   bool
}

type BookingStatus struct {
	BookingID string
	Status    Status
}

type Filter struct {
	TrainerID  *string // чтобы тренер видел свои тренировки
	BookedBy   *string
	BookedFrom *time.Time
	BookedTo   *time.Time
	WithLock   bool
}

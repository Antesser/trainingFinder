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

type Filter struct {
	BookedBy   string
	BookedFrom *time.Time
	BookedTo   *time.Time
	WithLock   bool
}

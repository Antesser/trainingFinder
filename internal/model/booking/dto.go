package booking

import "github.com/Antesser/trainingFinder/internal/model/page"

type ListBookingsRequest struct {
	Page     page.Page
	BookedBy string
}
type ListBookingsResponse struct {
	ModelList []TrainingBooking
	HasNext   bool
}

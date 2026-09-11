package booking

import (
	model "github.com/Antesser/trainingFinder/internal/model/booking"
)

type service struct {
	repo              bookingRepository
	outboxRepo        outboxRepository
	bookingMarshaller bookingMarshaller
}
type bookingMarshaller func(event model.CreateBookingEvent) ([]byte, error)

func New(bookingRepo bookingRepository, outboxRepo outboxRepository, bookingMarshaller bookingMarshaller) *service {
	return &service{repo: bookingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

package booking

import (
	model "github.com/Antesser/trainingFinder/internal/model/booking"
)

type service struct {
	repo              bookingRepository
	trainingRepo      trainingRepository
	outboxRepo        outboxRepository
	bookingMarshaller bookingMarshaller
}
type bookingMarshaller func(event model.CreateBookingEvent) ([]byte, error)

func New(bookingRepo bookingRepository, trainingRepo trainingRepository, bookingMarshaller bookingMarshaller, outboxRepo outboxRepository) *service {
	return &service{repo: bookingRepo, trainingRepo: trainingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

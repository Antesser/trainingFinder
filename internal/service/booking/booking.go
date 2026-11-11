package booking

import (
	"github.com/Antesser/trainingFinder/internal/model/booking"
)

type service struct {
	repo                bookingRepository
	trainingBookingRepo trainingBookingRepository
	outboxRepo          outboxRepository
	bookingMarshaller   marshalTrainingNotification
}
type marshalTrainingNotification func(model *booking.TrainingBookingStatus) ([]byte, error)

func New(bookingRepo bookingRepository, trainingBookingRepo trainingBookingRepository, bookingMarshaller marshalTrainingNotification, outboxRepo outboxRepository) *service {
	return &service{repo: bookingRepo, trainingBookingRepo: trainingBookingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

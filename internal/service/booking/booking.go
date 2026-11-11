package booking

import (
	"github.com/Antesser/trainingFinder/internal/model/booking"
)

type service struct {
	repo              bookingRepository
	trainingRepo      trainingRepository
	outboxRepo        outboxRepository
	bookingMarshaller marshalTrainingNotification
}
type marshalTrainingNotification func(model *booking.TrainingBookingStatus) ([]byte, error)

func New(bookingRepo bookingRepository, trainingRepo trainingRepository, bookingMarshaller marshalTrainingNotification, outboxRepo outboxRepository) *service {
	return &service{repo: bookingRepo, trainingRepo: trainingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

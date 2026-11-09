package booking

import (
	"github.com/Antesser/trainingFinder/internal/model/training"
)

type service struct {
	repo              bookingRepository
	trainingRepo      trainingRepository
	outboxRepo        outboxRepository
	bookingMarshaller MarshalTrainingNotification
}
type MarshalTrainingNotification func(model *training.Training) ([]byte, error)

func New(bookingRepo bookingRepository, trainingRepo trainingRepository, bookingMarshaller MarshalTrainingNotification, outboxRepo outboxRepository) *service {
	return &service{repo: bookingRepo, trainingRepo: trainingRepo, outboxRepo: outboxRepo, bookingMarshaller: bookingMarshaller}
}

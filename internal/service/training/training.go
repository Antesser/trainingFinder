package training

import (
	model "github.com/Antesser/trainingFinder/internal/model/training"
)

type trainingMarshaller func(event model.CreateTrainingEvent) ([]byte, error)

type trainingService struct {
	repo               trainingRepository
	outboxRepo         outboxRepository
	trainingMarshaller trainingMarshaller
}

func New(t trainingRepository, trainingMarshaller trainingMarshaller, outboxRepository outboxRepository) *trainingService {
	return &trainingService{repo: t, trainingMarshaller: trainingMarshaller, outboxRepo: outboxRepository}
}

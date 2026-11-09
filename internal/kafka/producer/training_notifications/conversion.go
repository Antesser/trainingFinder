package training_notifications

import (
	"fmt"

	"github.com/Antesser/trainingFinder/internal/model/training"
	kafkapb "github.com/Antesser/trainingFinder/pkg/api/kafka/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

func MarshalTrainingNotification(model *training.Training) ([]byte, error){

	msg := toProtoMessage(model)

	raw, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {
		return raw, fmt.Errorf("failed to marshal proto to json: %w", err)
	}

	return raw, nil
}

func toProtoMessage(training *training.Training) *kafkapb.TrainingData {
	return &kafkapb.TrainingData{
TrainingTime: training.Duration,
TrainingName: training.Name
	}
}

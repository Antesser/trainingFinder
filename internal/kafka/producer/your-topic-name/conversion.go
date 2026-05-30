package your_topic_name

import (
	"trainingFinder/internal/model/training"
	kafkapb "trainingFinder/pkg/api/kafka/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

func MarshalCreateTrainingEvent(event training.CreateTrainingEvent) ([]byte, error) {
	var msg *kafkapb.CreateTrainingEvent

	msg = &kafkapb.CreateTrainingEvent{
		TrainingId: event.TrainingID,
	}

	bytes, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {

	}

	return bytes, nil
}

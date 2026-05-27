package your_topic_name

import (
	kafkapb "pkg/api/kafka/v1"
	"trainingFinder/internal/model/training"

	"google.golang.org/protobuf/encoding/protojson"
)

func MarshalCreateTrainingEvent(event training.CreateTrainingEvent) ([]byte, error) {
	var msg *kafkapb.CreateTrainingEvent

	msg = &kafkfpb.CreateTrainingEvent{
		TrainigId: event.TrainingID,
	}

	bytes, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {

	}

	return bytes, nil
}

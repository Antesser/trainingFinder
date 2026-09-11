package your_topic_name

import (
	"github.com/Antesser/trainingFinder/internal/model/booking"
	"github.com/Antesser/trainingFinder/internal/model/training"
	kafkapb "github.com/Antesser/trainingFinder/pkg/api/kafka/v1"

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
func MarshalCreateBookingEvent(event booking.CreateBookingEvent) ([]byte, error) {
	var msg *kafkapb.CreateBookingEvent

	msg = &kafkapb.CreateBookingEvent{
		BookingId: event.BookingID,
	}

	bytes, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {

	}

	return bytes, nil
}

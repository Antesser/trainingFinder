package training_notifications

import (
	"fmt"

	modelBooking "github.com/Antesser/trainingFinder/internal/model/booking"
	kafkapb "github.com/Antesser/trainingFinder/pkg/api/kafka/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var statusConversion = map[modelBooking.Status]kafkapb.EventType{
	modelBooking.StatusUnspecified:         kafkapb.EventType_EVENT_TYPE_UNSPECIFIED,
	modelBooking.StatusAwaitingForApproval: kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_CREATED,
	modelBooking.StatusCancelled:           kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_CANCELLED,
	modelBooking.StatusApproved:            kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_APPROVED,
	modelBooking.StatusInProgress:          kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_IN_PROGRESS,
	modelBooking.StatusCompleted:           kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_COMPLETED,
	modelBooking.StatusDeclined:            kafkapb.EventType_EVENT_TYPE_TRAINING_BOOKING_DECLINED,
}

// описать мапу и в create\change_status вызов засунуть MarshalTrainingNotification и протестировать (сортировка uuid по времени created_at)

func MarshalTrainingNotification(model *modelBooking.TrainingBookingStatus) ([]byte, error) {

	msg := toProtoMessage(model)

	raw, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(msg)
	if err != nil {
		return raw, fmt.Errorf("failed to marshal proto to json: %w", err)
	}

	return raw, nil
}

func toProtoMessage(training *modelBooking.TrainingBookingStatus) *kafkapb.TrainingNotificationMessages {

	return &kafkapb.TrainingNotificationMessages{UserId: training.UserId, Destination: kafkapb.ContactType_CONTACT_TYPE_TELEGRAM, Data: &kafkapb.TrainingData{
		TrainingTime: training.Duration,
		TrainingName: training.TrainingName,
	}, NotificationType: statusConversion[training.Status]}
}

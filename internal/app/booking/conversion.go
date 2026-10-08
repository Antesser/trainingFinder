package booking

import (
	model "github.com/Antesser/trainingFinder/internal/model/booking"
	pb "github.com/Antesser/trainingFinder/pkg/api/booking/v1"
)

var (
	bookingStatusToModel = map[pb.BookingStatus]model.Status{
		pb.BookingStatus_BOOKING_STATUS_UNSPECIFIED:           model.StatusUnspecified,
		pb.BookingStatus_BOOKING_STATUS_AWAITING_FOR_APPROVAL: model.StatusAwaitingForApproval,
		pb.BookingStatus_BOOKING_STATUS_APPROVED:              model.StatusApproved,
		pb.BookingStatus_BOOKING_STATUS_IN_PROGRESS:           model.StatusInProgress,
		pb.BookingStatus_BOOKING_STATUS_COMPLETED:             model.StatusCompleted,
		pb.BookingStatus_BOOKING_STATUS_CANCELLED:             model.StatusCancelled,
		pb.BookingStatus_BOOKING_STATUS_DECLINED:              model.StatusDeclined,
	}
)

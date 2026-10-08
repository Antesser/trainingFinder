package booking

import "slices"

type (
	Status string
)

const (
	StatusUnspecified         Status = "UNSPECIFIED"
	StatusAwaitingForApproval Status = "AWAITING_FOR_APPROVAL"
	StatusCancelled           Status = "CANCELLED"
	StatusApproved            Status = "APPROVED"
	StatusInProgress          Status = "IN_PROGRESS"
	StatusCompleted           Status = "COMPLETED"
	StatusDeclined            Status = "DECLINED"
)

var (
	allowedStatusTransitions = map[Status][]Status{
		StatusAwaitingForApproval: {StatusApproved, StatusDeclined},
		StatusApproved:            {StatusInProgress, StatusCancelled},
		StatusInProgress:          {StatusCompleted},
		StatusCancelled:           {},
		StatusCompleted:           {},
		StatusDeclined:            {},
	}
)

func StatusTransitionAllowed(from, to Status) bool {
	return slices.Contains(allowedStatusTransitions[from], to)
}

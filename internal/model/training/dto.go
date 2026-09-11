package training

import (
	"time"

	"github.com/Antesser/trainingFinder/internal/model/page"
)

type UpdateTrainingRequest struct {
	ID             string
	TrainerID      *string
	UserID         *string
	StartedAt      *time.Time
	EndedAt        *time.Time
	AdditionalInfo *string
}

type trainingFilter struct {
	UserID   string
	Duration int
}
type ListTrainingRequest struct {
	Page   page.Page
	Filter trainingFilter
}
type ListTrainingResponse struct {
	ModelList []Training
	HasNext   bool
}

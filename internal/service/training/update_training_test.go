package training

import (
	"context"
	"strconv"
	"testing"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/training"

	"github.com/Antesser/trainingFinder/internal/service/training/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBehaviour func(m *mocks.TrainingRepository, model model.UpdateTrainingRequest)
type testCase struct {
	userModel  model.UpdateTrainingRequest
	resultUser error
	mock       mockBehaviour
}

func TestUpdateTraining(t *testing.T) {
	tt := []testCase{
		{userModel: model.UpdateTrainingRequest{ID: "jashdgfhj", TrainerID: new("jksadfhjk"), UserID: new("sjkdfhjk"), StartedAt: new(time.Now()), EndedAt: new(time.Now()), AdditionalInfo: new("some_info")},
			resultUser: nil,
			mock: func(m *mocks.TrainingRepository, model model.UpdateTrainingRequest) {
				m.On("UpdateTraining", mock.Anything, model).Return(nil)
			},
		}}
	for idx, testC := range tt {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			testRepo := mocks.NewTrainingRepository(t)
			testC.mock(testRepo, testC.userModel)
			svc := New(testRepo, nil, nil)
			err := svc.UpdateTraining(context.Background(), testC.userModel)
			assert.NoError(t, err)
		})
	}

}

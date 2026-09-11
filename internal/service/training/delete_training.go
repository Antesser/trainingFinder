package training

import "context"

func (t *trainingService) DeleteTraining(ctx context.Context, id string) error {
	err := t.repo.DeleteTraining(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

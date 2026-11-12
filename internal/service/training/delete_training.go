package training

import (
	"context"
	"fmt"
)

func (t *trainingService) DeleteTraining(ctx context.Context, id string) error {
	if err := t.repo.DeleteTraining(ctx, id); err != nil {
		return fmt.Errorf("delete training: %w", err)
	}
	return nil
}

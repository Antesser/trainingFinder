package users

import (
	"context"
	model "trainingFinder/internal/model/user"
)

type userService interface {
	GetUserByID(ctx context.Context, id string) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id, login string) error
}

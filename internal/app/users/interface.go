package users

import "context"

type userService interface {
	GetUser(ctx context.Context, id string) (string, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id, username string) error
}

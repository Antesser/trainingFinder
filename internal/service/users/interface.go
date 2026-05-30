package users

import "context"

type userRepository interface {
	GetUserByID(ctx context.Context, id string) (string, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id, username string) error
}

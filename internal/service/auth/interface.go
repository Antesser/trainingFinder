package auth

import "context"

type authRepository interface {
	SignUp(ctx context.Context, login, password string) (string, string, error)
}

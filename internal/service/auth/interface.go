package auth

import "context"

type authRepository interface {
	SignUp(ctx context.Context, hash, login, password string) (string, string, error)
}

package auth

import "context"

type authService interface {
	SignUp(ctx context.Context, login, password string) (string, error)
	SignIn(ctx context.Context, login, password string) (string, string, error)
}

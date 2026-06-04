package auth

import (
	"context"
	model "trainingFinder/internal/model/auth"
)

type authService interface {
	SignUp(ctx context.Context, login, password string) (string, error)
	SignIn(ctx context.Context, login, password string) (model.Tokens, error)
}

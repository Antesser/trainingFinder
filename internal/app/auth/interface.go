package auth

import (
	"context"
	model "trainingFinder/internal/model/auth"

	"github.com/google/uuid"
)

type authService interface {
	SignUp(ctx context.Context, login, password string) (string, error)
	SignIn(ctx context.Context, login, password string) (model.Tokens, error)
	RefreshSession(ctx context.Context, refreshToken uuid.UUID) (accessToken string, err error)
}

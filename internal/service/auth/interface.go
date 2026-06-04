package auth

import (
	"context"
	model "trainingFinder/internal/model/auth"
)

type authRepository interface {
	SignUp(ctx context.Context, hash []byte, id, login string) (string, error)
	GetUserAuthInfoByLogin(ctx context.Context, login string) (model.UserAuthInfo, error)
}

package auth

import (
	"context"
	model "trainingFinder/internal/model/auth"

	"github.com/google/uuid"
)

type authRepository interface {
	SignUp(ctx context.Context, hash []byte, id, login string) (string, error)
	GetUserAuthInfoByLogin(ctx context.Context, login string) (model.UserAuthInfo, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (*model.Sessions, error)
	SaveRefreshToken(ctx context.Context, session model.Sessions) error
}

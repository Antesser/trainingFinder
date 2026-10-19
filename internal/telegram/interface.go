package telegram

import (
	"context"

	redisModel "github.com/Antesser/trainingFinder/internal/model/telegram"
)

type authService interface {
	SignUp(ctx context.Context, login, password, roleID string) (string, error)
}
type sessionCache interface {
	SaveSession(ctx context.Context, chatID int64, s *redisModel.Session) error
	LoadSession(ctx context.Context, chatID int64) (*redisModel.Session, error)
	DelSession(ctx context.Context, chatID int64) error
}

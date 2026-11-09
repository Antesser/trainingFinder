package auth

import (
	"context"

	"github.com/google/uuid"

	model "github.com/Antesser/trainingFinder/internal/model/auth"
	"github.com/Antesser/trainingFinder/internal/utils"
)

func (s *service) RefreshSession(ctx context.Context, refreshToken uuid.UUID) (accessToken string, err error) {
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	if !session.IsActive() {
		return "", model.ErrSessionExpired
	}
	aToken, err := utils.GenerateToken(session.UserID, session.RoleID, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return aToken, nil
}

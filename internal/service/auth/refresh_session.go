package auth

import (
	"context"
	"fmt"

	model "github.com/Antesser/trainingFinder/internal/model/auth"
	"github.com/Antesser/trainingFinder/internal/utils"
	"github.com/google/uuid"
)

func (s *service) RefreshSession(ctx context.Context, refreshToken uuid.UUID) (string, error) {
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("get session by refresh token: %w", err)
	}
	if !session.IsActive() {
		return "", model.ErrSessionExpired
	}
	aToken, err := utils.GenerateToken(session.UserID, session.RoleID, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return aToken, nil
}

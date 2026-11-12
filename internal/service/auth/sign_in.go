package auth

import (
	"context"
	"fmt"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/auth"
	"github.com/Antesser/trainingFinder/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignIn(ctx context.Context, login, password string) (model.Tokens, error) {
	authInfo, err := s.authRepo.GetUserAuthInfoByLogin(ctx, login)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("get user auth info by login: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(authInfo.Password), []byte(password)); err != nil {
		return model.Tokens{}, model.ErrInvalidPass
	}

	aToken, err := utils.GenerateToken(authInfo.ID, authInfo.RoleID, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return model.Tokens{}, fmt.Errorf("generate access token: %w", err)
	}

	refToken := uuid.New()
	refreshTokenModel := model.Session{
		Token:     refToken,
		UserID:    authInfo.ID,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().Add(s.accessTokenDuration),
	}
	if err = s.authRepo.CreateSession(ctx, refreshTokenModel); err != nil {
		return model.Tokens{}, fmt.Errorf("create session: %w", err)
	}

	return model.Tokens{AccessToken: aToken, RefreshToken: refToken.String()}, nil
}

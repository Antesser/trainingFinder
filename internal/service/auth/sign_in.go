package auth

import (
	"context"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/auth"
	"github.com/Antesser/trainingFinder/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignIn(ctx context.Context, login, password string) (model.Tokens, error) {
	authInfo, err := s.authRepo.GetUserAuthInfoByLogin(ctx, login)
	if err != nil {
		return model.Tokens{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(authInfo.Password), []byte(password))
	if err != nil {
		return model.Tokens{}, model.ErrInvalidPass
	}
	aToken, err := utils.GenerateToken(authInfo.ID, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return model.Tokens{}, err
	}
	refToken := uuid.New()
	createdAt := time.Now().UTC()
	refreshTokenModel := model.Session{Token: refToken, UserID: authInfo.ID, Active: true, CreatedAt: createdAt, ExpiresAt: time.Now().Add(s.accessTokenDuration)}
	err = s.authRepo.CreateSession(ctx, refreshTokenModel)
	if err != nil {
		return model.Tokens{}, err
	}
	return model.Tokens{AccessToken: aToken, RefreshToken: refToken.String()}, nil
}

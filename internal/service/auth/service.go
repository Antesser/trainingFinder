package auth

import (
	"context"
	"time"
	model "trainingFinder/internal/model/auth"
	"trainingFinder/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	authRepo            authRepository
	secretKey           string
	accessTokenDuration time.Duration
}

func New(authRepo authRepository, skey string, atd time.Duration) *service {
	return &service{
		authRepo:            authRepo,
		accessTokenDuration: atd,
		secretKey:           skey,
	}
}

func (s *service) SignIn(ctx context.Context, login, password string) (model.Tokens, error) {
	// создать модельку userauthinfo (внутри репозитория) и работать с ней, проверять хешированный пароль из базы с переданным
	// refresh token через uuid сделать отдельной таблицей в БД по гайду от негодующего Виталия из ТГ
	authInfo, err := s.authRepo.GetUserAuthInfoByLogin(ctx, login) // пароль захешировать прям тута, Виталий снова негодует
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
	// refToken создать токен, модельку сессии, сохранить сессию в БД и вернуть оба токена пользователю
	refToken := uuid.New()
	refreshTokenModel := model.Sessions{Token: refToken, UserID: authInfo.ID, Active: true, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(s.accessTokenDuration)}
	err = s.authRepo.SaveRefreshToken(ctx, refreshTokenModel) // пароль захешировать прям тута, Виталий снова негодует
	if err != nil {
		return model.Tokens{}, err
	}
	return model.Tokens{AccessToken: aToken, RefreshToken: refToken.String()}, nil
}
func (s *service) SignUp(ctx context.Context, login, password string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userId := uuid.New().String()
	id, err := s.authRepo.SignUp(ctx, hash, userId, login) // пароль захешировать, Виталий снова негодует
	if err != nil {
		return "", err
	}
	//aToken, err := utils.GenerateToken(hash, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (s *service) RefreshSession(ctx context.Context, refreshToken uuid.UUID) (accessToken string, err error) {
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	if !session.IsActive() {
		return "", model.ErrSessionExpired
	}
	aToken, err := utils.GenerateToken(session.UserID, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return aToken, nil
}

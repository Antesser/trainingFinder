package auth

import (
	"context"
	"time"

	model "github.com/Antesser/trainingFinder/internal/model/auth"
	"github.com/Antesser/trainingFinder/internal/utils"

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
func (s *service) SignUp(ctx context.Context, login, password string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	userId := uuid.New().String()
	id, err := s.authRepo.SignUp(ctx, hash, userId, login)
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

func (s *service) CreateRole(ctx context.Context, role string) error {
	err := s.authRepo.CreateRole(ctx, role)
	if err != nil {
		return err
	}
	return nil
}

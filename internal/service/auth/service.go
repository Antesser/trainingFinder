package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
	model "trainingFinder/internal/model/auth"
	"trainingFinder/internal/utils"

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
func generateRandomHash(n int) (string, error) { // перенести в utils/helper (pkg, если переиспользоваться в других сервисах)
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *service) SignIn(ctx context.Context, login, password string) (model.Tokens, error) {
	// создать модельку userauthinfo (внутри репозитория) и работать с ней, проверять хешированный пароль из базы с переданным
	// refresh token через uuid сделать отдельной таблицей в БД по гайду от негодующего Виталия из ТГ
	authInfo, err := s.authRepo.GetUserByIDAuthInfo(ctx, login) // пароль захешировать прям тута, Виталий снова негодует
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
	return model.Tokens{AccessToken: aToken}, nil
}
func (s *service) SignUp(ctx context.Context, login, password string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	id, err := s.authRepo.SignUp(ctx, hash, login, password) // пароль захешировать, Виталий снова негодует
	if err != nil {
		return "", err
	}
	//aToken, err := utils.GenerateToken(hash, []byte(s.secretKey), s.accessTokenDuration)
	if err != nil {
		return "", err
	}
	return id, nil
}

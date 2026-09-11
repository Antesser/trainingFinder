package auth

import (
	"time"
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

package auth

import (
	"errors"
)

type UserAuthInfo struct {
	ID       string
	Login    string
	Password string
	RoleID   string
}

var (
	ErrNotFound    = errors.New("user not found")
	ErrInvalidPass = errors.New("password is invalid")
)

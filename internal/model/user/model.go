package user

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	Id       string
	Username string
}

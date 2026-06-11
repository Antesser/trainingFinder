package user

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrLoginAlreadyExists = errors.New("login already exists")

type User struct {
	Id    string
	Login string
}

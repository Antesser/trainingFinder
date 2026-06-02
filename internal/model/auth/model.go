package auth

import (
	"fmt"
)

type UserAuthInfo struct {
	ID       string
	Login    string
	Password string
}

var ErrNotFound error = fmt.Errorf("user not found")
var ErrInvalidPass error = fmt.Errorf("password is invalid")

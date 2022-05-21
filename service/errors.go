package service

import (
	"errors"
)

var (
	ErrPasswordIncorrect   = errors.New("Password incorrect")
	ErrUserAlreadyRegister = errors.New("User already register")
)

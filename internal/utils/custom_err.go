package utils

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserAlreadyLoggedOut = errors.New("user already logged Out")
	ErrUserNotFound         = errors.New("user not found")
)

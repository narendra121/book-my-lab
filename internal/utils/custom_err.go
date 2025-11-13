package utils

import "errors"

var (
	// ErrUserAlreadyExists          = errors.New("user already exists")
	ErrUserAlreadyLoggedOut               = errors.New("user already logged Out")
	ErrUserNotFound                       = errors.New("user not found")
	ErrUserAlreadyDeleted                 = errors.New("user already deleted")
	ErrInvalidUserOrPass                  = errors.New("invalid username or password")
	ErrUserAlreadyExistsWithEmail         = errors.New("user already exists with email id")
	ErrUserAlreadyExistsWithPhone         = errors.New("user already exists with phone number")
	ErrUserAlreadyExistsWithEmailAndPhone = errors.New("user already exists with email id and phone number")
	ErrUserAlreadyExistsWithEmailOrPhone  = errors.New("user already exists with email id or phone number")
)

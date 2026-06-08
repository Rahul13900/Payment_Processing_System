package usererrors

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be atleast 8 characters")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserSuspended      = errors.New("user account is suspended")
	ErrInternalError      = errors.New("internal server error")
)

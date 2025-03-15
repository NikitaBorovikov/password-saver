package apperrors

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDatabaseInternal = errors.New("internal database error")
	ErrDuplicateUser    = errors.New("user already exists")
	ErrValidateUser     = errors.New("failed to validate user data")
	ErrHashPassword     = errors.New("failed to hash password")
	ErrComparePasswords = errors.New("incorrect password")
	ErrDecodeRequest    = errors.New("failed to decode request")
	ErrServerInternal   = errors.New("internal servr error")
	ErrNotAuthenticated = errors.New("user is not authenticated")
	ErrInvalidURLParam  = errors.New("invalid URL param")

	ErrPasswordNotExists = errors.New("password not exists")
	ErrValidatePassword  = errors.New("failed to validate password")
)

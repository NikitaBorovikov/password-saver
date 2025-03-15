package apperrors

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrDatabaseInternal = errors.New("internal database error")
	ErrDuplicateUser    = errors.New("user already exists")
	ErrValidateUser     = errors.New("failed to validate user data")
	ErrHashPassword     = errors.New("failed to hash password")
	ErrComparePasswords = errors.New("incorrect password")
)

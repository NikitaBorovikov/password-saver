package apperrors

import "errors"

// Database errors
var (
	ErrUserNotFound     = errors.New("failed to find user")
	ErrDuplicateUser    = errors.New("user already exists")
	ErrDatabaseInternal = errors.New("internal database error")
)

// Validation errors
var (
	ErrValidateUser     = errors.New("failed to validate user data")
	ErrValidatePassword = errors.New("failed to validate password")
)

// Auth errors
var (
	ErrNotAuthenticated = errors.New("user is not authenticated")
	ErrHashPassword     = errors.New("failed to hash password")
	ErrComparePasswords = errors.New("incorrect password")
)

// Request errors
var (
	ErrDecodeRequest   = errors.New("failed to decode request")
	ErrInvalidURLParam = errors.New("invalid URL param")
)

// Password errors
var (
	ErrPasswordNotExists = errors.New("password not exists")
)

// Server errors
var (
	ErrServerInternal = errors.New("internal server error")
)

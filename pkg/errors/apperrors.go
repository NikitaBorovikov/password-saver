package apperrors

import "errors"

// Database errors
var (
	ErrUserNotFound     = errors.New("failed to find user")
	ErrDuplicateUser    = errors.New("user already exists")
	ErrDatabaseInternal = errors.New("internal database error")
)

// User validation errors
var (
	ErrValidateUser              = errors.New("failed to validate user data")
	ErrValidateEmailField        = errors.New("failed to validate email: the email field is required and must be an email address")
	ErrValidateUserPasswordField = errors.New("failed to validate user password: password field is required and must be between 7 and 40 characters long")
	ErrValidateOldPasswordField  = errors.New("failed to validate old password: the old password field is required")
	ErrValidateNewPasswordField  = errors.New("failed to validate new password: the new password field is required and must be between 7 and 40 characters long")
)

// Password validation errors
var (
	ErrValidatePassword          = errors.New("failed to validate password data")
	ErrValidateServiceField      = errors.New("failed to validate service: the service field must be between 1 and 100 characters long")
	ErrValidateSavePasswordField = errors.New("failed to validate password: the password field must be between 1 and 100 characters long")
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

package usecases

import (
	"errors"
	"password-saver/internal/infrastructure/repository/postgres"
)

var (
	// Query errors
	ErrNotFound      = errors.New("failed to find data")
	ErrAlreadyExists = errors.New("data already exists")
	// Data processing errors
	ErrDataProcessing       = errors.New("data processing error")
	ErrHashPassword         = errors.New("failed to hash password")
	ErrEcryptData           = errors.New("failed to encrypt data")
	ErrDecryptData          = errors.New("failed to decrypt data")
	ErrMakePasswordResponse = errors.New("failed to make password response")
	// Validation errors
	ErrInvalidInput = errors.New("invalid input data")
	// Auth errors
	ErrComparePasswords = errors.New("wrong password")
	// System errors
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrInternalDB         = errors.New("internal database error")
	// Service errors
	ErrPingDB = errors.New("failed to ping DB")
)

func handleRepositoryErrors(err error) error {
	switch {
	case errors.Is(err, postgres.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, postgres.ErrDuplicateData):
		return ErrAlreadyExists
	case errors.Is(err, postgres.ErrScanFailed), errors.Is(err, postgres.ErrForeignKeyViolation):
		return ErrDataProcessing
	case errors.Is(err, postgres.ErrConnectionFailed):
		return ErrServiceUnavailable
	case errors.Is(err, postgres.ErrPingFailed):
		return ErrPingDB
	default:
		return ErrInternalDB
	}
}

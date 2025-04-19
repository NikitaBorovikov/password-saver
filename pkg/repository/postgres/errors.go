package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
)

const (
	uniqueViolationErrCode  = "23505"
	pingFailedCode          = "08006"
	foreignKeyViolationCode = "23503"
)

// Custom repository errors
var (
	ErrUserNotFound        = errors.New("failed to find user")
	ErrDuplicateUser       = errors.New("user already exists")
	ErrConnectionFailed    = errors.New("database connection failed")
	ErrPingFailed          = errors.New("failed to ping database")
	ErrForeignKeyViolation = errors.New("referenced entity not found")
	ErrScanFailed          = errors.New("failed to scan data")
	ErrInternalDB          = errors.New("internal database error")
)

// Handles SQL errors and returns custom repository errors
func handleSQLErrors(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	if strings.Contains(err.Error(), "scan") {
		return ErrScanFailed
	}

	// Handle specific PostgreSQL errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case uniqueViolationErrCode:
			return ErrDuplicateUser
		case pingFailedCode:
			return ErrConnectionFailed
		case foreignKeyViolationCode:
			return ErrForeignKeyViolation
		}
	}

	return ErrInternalDB
}

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
	// Query errors
	ErrNotFound            = errors.New("failed to find data")
	ErrDuplicateData       = errors.New("data already exists")
	ErrForeignKeyViolation = errors.New("referenced entity not found")
	// Connection errors
	ErrConnectionFailed = errors.New("database connection failed")
	ErrPingFailed       = errors.New("failed to ping database")
	// Data pocessing errors
	ErrScanFailed = errors.New("failed to scan data")

	ErrInternalDB = errors.New("internal database error")
)

// Handles SQL errors and returns custom repository errors
func handleSQLErrors(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	if strings.Contains(err.Error(), "scan") {
		return ErrScanFailed
	}

	// Handle specific PostgreSQL errors
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case uniqueViolationErrCode:
			return ErrDuplicateData
		case pingFailedCode:
			return ErrConnectionFailed
		case foreignKeyViolationCode:
			return ErrForeignKeyViolation
		}
	}

	return ErrInternalDB
}

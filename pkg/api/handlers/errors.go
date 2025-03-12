package handlers

import (
	"errors"
)

var (
	errUserIDNotInContext = errors.New("userID not found in request context")
	errStrToIntConversion = errors.New("failed to convert userID from string to int")
)

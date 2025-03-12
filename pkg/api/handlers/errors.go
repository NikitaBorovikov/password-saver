package handlers

import (
	"errors"
)

var (
	errUserIDNotInContext   = errors.New("userID not found in request context")
	errStrToIntConversion   = errors.New("failed to convert userID from string to int")
	errorNotAuthenticated   = errors.New("user is not authenticated")
	errorUserIDNotInSession = errors.New("user ID in session not found or invalid")
)

package handlers

import (
	"errors"
	"net/http"
	"password-saver/pkg/usecases"
)

var (
	ErrDecodeRequest    = errors.New("failed to decode request")
	ErrNotAuthenticated = errors.New("user is not authenticated")

	ErrorNotFound         = errors.New("requested data not found")
	ErrAlreadyExists      = errors.New("data already exists")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrWrongPassword      = errors.New("wrong password or email")
	ErrDataProcessing     = errors.New("data processing error")
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrInternalServer     = errors.New("internal server error")
)

func handleUsecaseErrors(err error) (int, error) {
	switch {
	case errors.Is(err, usecases.ErrNotFound):
		return http.StatusNotFound, ErrorNotFound
	case errors.Is(err, usecases.ErrAlreadyExists):
		return http.StatusConflict, ErrAlreadyExists
	case errors.Is(err, usecases.ErrHashPassword),
		errors.Is(err, usecases.ErrEcryptData),
		errors.Is(err, usecases.ErrDecryptData),
		errors.Is(err, usecases.ErrMakePasswordResponse):
		return http.StatusInternalServerError, ErrDataProcessing
	case errors.Is(err, usecases.ErrInvalidInput):
		return http.StatusUnprocessableEntity, ErrInvalidInput
	case errors.Is(err, usecases.ErrComparePasswords):
		return http.StatusUnauthorized, ErrWrongPassword
	default:
		return http.StatusInternalServerError, ErrInternalServer
	}
}
